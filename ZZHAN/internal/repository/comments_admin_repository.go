package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"
	"time"

	"gorm.io/gorm"
)

type commentsAdminRepository struct {
	db *gorm.DB
}

// NewCommentsAdminRepository 创建后台评论仓储实例
func NewCommentsAdminRepository(db *gorm.DB) CommentsAdminRepository {
	return &commentsAdminRepository{db: db}
}

// AdminList 获取评论列表
func (r *commentsAdminRepository) AdminList(ctx context.Context, req *dto.CommentsAdminListRequest) ([]dto.CommentsAdminItem, int64, error) {
	var comments []entity.Comment
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Comment{})

	// 关键词搜索（内容/用户名）
	if req.Keyword != "" {
		query = query.Where("(content LIKE ? OR user_name LIKE ?)", "%"+req.Keyword+"%", "%"+req.Keyword+"%")
	}

	// 按状态筛选
	if req.Status != "" && req.Status != "all" {
		query = query.Where("status = ?", req.Status)
	}

	// 按文章筛选
	if req.ArticleID > 0 {
		query = query.Where("article_id = ?", req.ArticleID)
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
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	result := make([]dto.CommentsAdminItem, 0, len(comments))
	for _, c := range comments {
		item := dto.CommentsAdminItem{
			ID:         c.ID,
			ArticleID:  c.ArticleID,
			ParentID:   c.ParentID,
			UserID:     c.UserID,
			UserName:   c.UserName,
			UserAvatar: c.UserAvatar,
			Content:    c.Content,
			Status:     c.Status,
			LikeCount:  c.LikeCount,
			IP:         c.IP,
			CreatedAt:  c.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		// 查询所属文章标题
		var article entity.Article
		if err := r.db.WithContext(ctx).Select("title").First(&article, c.ArticleID).Error; err == nil {
			item.ArticleTitle = article.Title
		}

		result = append(result, item)
	}

	return result, total, nil
}

// AdminGetByID 通过 ID 获取评论详情
func (r *commentsAdminRepository) AdminGetByID(ctx context.Context, id int64) (*dto.CommentsAdminItem, error) {
	var comment entity.Comment

	if err := r.db.WithContext(ctx).First(&comment, id).Error; err != nil {
		return nil, err
	}

	item := &dto.CommentsAdminItem{
		ID:         comment.ID,
		ArticleID:  comment.ArticleID,
		ParentID:   comment.ParentID,
		UserID:     comment.UserID,
		UserName:   comment.UserName,
		UserAvatar: comment.UserAvatar,
		Content:    comment.Content,
		Status:     comment.Status,
		LikeCount:  comment.LikeCount,
		IP:         comment.IP,
		CreatedAt:  comment.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	// 查询所属文章标题
	var article entity.Article
	if err := r.db.WithContext(ctx).Select("title").First(&article, comment.ArticleID).Error; err == nil {
		item.ArticleTitle = article.Title
	}

	return item, nil
}

// AdminUpdateStatus 修改评论状态
func (r *commentsAdminRepository) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Comment{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// AdminDelete 删除评论（软删除）
func (r *commentsAdminRepository) AdminDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Comment{}, id).Error
}
