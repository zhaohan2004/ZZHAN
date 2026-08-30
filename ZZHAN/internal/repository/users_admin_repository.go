package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"
	"time"

	"gorm.io/gorm"
)

type usersAdminRepository struct {
	db *gorm.DB
}

// NewUsersAdminRepository 创建后台用户仓储实例
func NewUsersAdminRepository(db *gorm.DB) UsersAdminRepository {
	return &usersAdminRepository{db: db}
}

// AdminList 获取用户列表
func (r *usersAdminRepository) AdminList(ctx context.Context, req *dto.UsersAdminListRequest) ([]dto.UsersAdminItem, int64, error) {
	var users []entity.User
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.User{})

	// 关键词搜索（昵称）
	if req.Keyword != "" {
		query = query.Where("nickname LIKE ?", "%"+req.Keyword+"%")
	}

	// 按状态筛选
	if req.Status != "" && req.Status != "all" {
		if req.Status == "active" {
			query = query.Where("status = ?", 1)
		} else if req.Status == "inactive" {
			query = query.Where("status = ?", 0)
		}
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
		Find(&users).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	result := make([]dto.UsersAdminItem, 0, len(users))
	for _, u := range users {
		item := dto.UsersAdminItem{
			ID:        u.ID,
			Provider:  u.Provider,
			Openid:    u.Openid,
			Nickname:  u.Nickname,
			Avatar:    u.Avatar,
			Status:    u.Status,
			CreatedAt: u.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		if u.LastLoginAt != nil {
			item.LastLoginAt = u.LastLoginAt.Format("2006-01-02 15:04:05")
		}

		result = append(result, item)
	}

	return result, total, nil
}

// AdminGetByID 通过 ID 获取用户详情
func (r *usersAdminRepository) AdminGetByID(ctx context.Context, id int64) (*dto.UsersAdminItem, error) {
	var user entity.User

	if err := r.db.WithContext(ctx).First(&user, id).Error; err != nil {
		return nil, err
	}

	item := &dto.UsersAdminItem{
		ID:        user.ID,
		Provider:  user.Provider,
		Openid:    user.Openid,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Status:    user.Status,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if user.LastLoginAt != nil {
		item.LastLoginAt = user.LastLoginAt.Format("2006-01-02 15:04:05")
	}

	return item, nil
}

// AdminUpdateStatus 修改用户状态
func (r *usersAdminRepository) AdminUpdateStatus(ctx context.Context, id int64, status int8) error {
	return r.db.WithContext(ctx).Model(&entity.User{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// AdminDelete 删除用户
func (r *usersAdminRepository) AdminDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.User{}, id).Error
}
