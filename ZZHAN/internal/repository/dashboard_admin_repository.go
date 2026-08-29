package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"

	"gorm.io/gorm"
)

type dashboardAdminRepository struct {
	db *gorm.DB
}

// NewDashboardAdminRepository 创建后台仪表盘仓储实例
func NewDashboardAdminRepository(db *gorm.DB) DashboardAdminRepository {
	return &dashboardAdminRepository{db: db}
}

// GetStats 获取统计数据
func (r *dashboardAdminRepository) GetStats(ctx context.Context) (map[string]dto.StatItem, error) {
	result := make(map[string]dto.StatItem)

	// --- 文章总数 ---
	var articleCount int64
	r.db.WithContext(ctx).Model(&entity.Article{}).Count(&articleCount)
	result["articles"] = dto.StatItem{Value: articleCount}

	// --- 分类数量 ---
	var categoryCount int64
	r.db.WithContext(ctx).Model(&entity.Category{}).Where("status = ?", "active").Count(&categoryCount)
	result["categories"] = dto.StatItem{Value: categoryCount}

	// --- 标签数量 ---
	var tagCount int64
	r.db.WithContext(ctx).Model(&entity.Tag{}).Where("status = ?", "active").Count(&tagCount)
	result["tags"] = dto.StatItem{Value: tagCount}

	// --- 评论数量 ---
	var commentCount int64
	r.db.WithContext(ctx).Model(&entity.Comment{}).Count(&commentCount)
	result["comments"] = dto.StatItem{Value: commentCount}

	return result, nil
}

// GetArticles 获取最新发布文章
func (r *dashboardAdminRepository) GetArticles(ctx context.Context) ([]dto.RecentPostItem, error) {
	// 最新5篇已发布文章
	var articles []entity.Article
	if err := r.db.WithContext(ctx).
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&articles).Error; err != nil {
		return nil, err
	}

	posts := make([]dto.RecentPostItem, 0, len(articles))
	for _, a := range articles {
		item := dto.RecentPostItem{
			ID:    a.ID,
			Title: a.Title,
			Date:  a.PublishedAt.Format("2006-01-02"),
			Views: a.Views,
		}
		// 查询分类名
		var cat entity.Category
		if err := r.db.WithContext(ctx).Select("name").First(&cat, a.CategoryID).Error; err == nil {
			item.Category = cat.Name
		}
		posts = append(posts, item)
	}

	return posts, nil
}

// GetComments 获取最新5条评论
func (r *dashboardAdminRepository) GetComments(ctx context.Context) ([]dto.DashboardCommentItem, error) {
	var comments []entity.Comment
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(5).
		Find(&comments).Error; err != nil {
		return nil, err
	}

	items := make([]dto.DashboardCommentItem, 0, len(comments))
	for _, c := range comments {
		items = append(items, dto.DashboardCommentItem{
			ID:       c.ID,
			UserName: c.UserName,
			Avatar:   c.UserAvatar,
			Content:  c.Content,
			Time:     c.CreatedAt.Format("2006-01-02 15:04"),
		})
	}

	return items, nil
}

// GetOperations 获取最新8条操作记录
func (r *dashboardAdminRepository) GetOperations(ctx context.Context) ([]dto.DashboardOperationItem, error) {
	var logs []entity.OperationLog
	if err := r.db.WithContext(ctx).
		Order("created_at DESC").
		Limit(8).
		Find(&logs).Error; err != nil {
		return nil, err
	}

	items := make([]dto.DashboardOperationItem, 0, len(logs))
	for _, l := range logs {
		item := dto.DashboardOperationItem{
			Time:   l.CreatedAt.Format("2006-01-02 15:04"),
			Action: l.Action,
			Target: l.Target,
		}
		// 查询管理员昵称
		if l.AdminID != nil {
			var admin entity.Admin
			if err := r.db.WithContext(ctx).Select("nickname").First(&admin, *l.AdminID).Error; err == nil {
				item.User = admin.Nickname
			}
		}
		items = append(items, item)
	}

	return items, nil
}
