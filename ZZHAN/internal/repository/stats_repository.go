package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"
	"sort"

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

	// 获取最近动态：最新文章
	var recentArticles []entity.Article
	r.db.WithContext(ctx).
		Where("status = ?", "published").
		Order("published_at DESC").
		Limit(5).
		Find(&recentArticles)

	dynamics := make([]dto.Dynamic, 0, len(recentArticles))
	for _, a := range recentArticles {
		if a.PublishedAt == nil {
			continue
		}
		dynamics = append(dynamics, dto.Dynamic{
			Type: "article",
			Text: a.Title,
			Time: a.PublishedAt.Format("2006-01-02 15:04:05"),
			Link: "/article/" + a.Slug,
		})
	}

	// 获取最近评论
	var recentComments []entity.Comment
	r.db.WithContext(ctx).
		Where("status = ?", "normal").
		Order("created_at DESC").
		Limit(5).
		Find(&recentComments)

	for _, c := range recentComments {
		dynamics = append(dynamics, dto.Dynamic{
			Type: "comment",
			Text: c.UserName + " 评论了文章",
			Time: c.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// 按时间倒序排列，取最近10条
	sort.Slice(dynamics, func(i, j int) bool {
		return dynamics[i].Time > dynamics[j].Time
	})
	if len(dynamics) > 10 {
		dynamics = dynamics[:10]
	}
	resp.Dynamics = dynamics

	return resp, nil
}
