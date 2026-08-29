package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"
	"time"

	"gorm.io/gorm"
)

type operationLogsAdminRepository struct {
	db *gorm.DB
}

// NewOperationLogsAdminRepository 创建后台操作日志仓储实例
func NewOperationLogsAdminRepository(db *gorm.DB) OperationLogsAdminRepository {
	return &operationLogsAdminRepository{db: db}
}

// AdminList 获取操作日志列表
func (r *operationLogsAdminRepository) AdminList(ctx context.Context, req *dto.OperationLogsAdminListRequest) ([]dto.OperationLogsAdminItem, int64, error) {
	var logs []entity.OperationLog
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.OperationLog{})

	// 按操作类型筛选
	if req.Action != "" {
		query = query.Where("action = ?", req.Action)
	}

	// 按操作对象模糊搜索
	if req.Target != "" {
		query = query.Where("target LIKE ?", "%"+req.Target+"%")
	}

	// 时间段筛选
	if req.StartDate != "" {
		if t, err := time.Parse("2006-01-02", req.StartDate); err == nil {
			query = query.Where("created_at >= ?", t)
		}
	}
	if req.EndDate != "" {
		if t, err := time.Parse("2006-01-02", req.EndDate); err == nil {
			// 包含当天，加一天减一秒
			query = query.Where("created_at < ?", t.Add(24*time.Hour))
		}
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	result := make([]dto.OperationLogsAdminItem, 0, len(logs))
	for _, l := range logs {
		item := dto.OperationLogsAdminItem{
			ID:        l.ID,
			AdminID:   l.AdminID,
			Action:    l.Action,
			Target:    l.Target,
			CreatedAt: l.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// 查询管理员昵称
		if l.AdminID != nil {
			var admin entity.Admin
			if err := r.db.WithContext(ctx).Select("nickname").First(&admin, *l.AdminID).Error; err == nil {
				item.AdminName = admin.Nickname
			}
		}

		result = append(result, item)
	}

	return result, total, nil
}
