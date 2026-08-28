package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"

	"gorm.io/gorm"
)

type statsRepository struct {
	db *gorm.DB
}

// NewStatsRepository 创建统计仓储实例
func NewStatsRepository(db *gorm.DB) StatsRepository {
	return &statsRepository{db: db}
}

// GetStats 获取站点统计数据
func (r *statsRepository) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	resp := &dto.StatsResponse{
		Dynamics: []dto.Dynamic{},
	}

	// 统计已发布文章总数
	var articleCount int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Article{}).
		Where("status = ?", "published").
		Count(&articleCount).Error; err != nil {
		return nil, err
	}
	resp.Articles = int(articleCount)

	// 统计总浏览量
	var totalViews int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Article{}).
		Where("status = ?", "published").
		Select("COALESCE(SUM(views), 0)").
		Scan(&totalViews).Error; err != nil {
		return nil, err
	}
	resp.Views = totalViews

	// 统计总评论数
	var totalComments int64
	if err := r.db.WithContext(ctx).
		Model(&entity.Comment{}).
		Where("status = ?", "normal").
		Count(&totalComments).Error; err != nil {
		return nil, err
	}
	resp.Comments = totalComments

	// 获取最近动态：最新文章
	var recentArticles []entity.Article
	r.db.WithContext(ctx).
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(8).
		Find(&recentArticles)

	dynamics := make([]dto.Dynamic, 0, len(recentArticles))
	for _, a := range recentArticles {
		if a.PublishedAt == nil {
			continue
		}
		dynamics = append(dynamics, dto.Dynamic{
			Type: "write",
			Text: "发布了文章《" + a.Title + "》",
			Time: a.PublishedAt.Format("2006-01-02"),
			Link: "/article/" + a.Slug,
		})
	}

	resp.Dynamics = dynamics

	return resp, nil
}
