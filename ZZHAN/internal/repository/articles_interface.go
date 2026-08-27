package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// ArticlesRepository 文章仓储接口
type ArticlesRepository interface {
	// GetPublishedList 获取已发布文章列表（分页、筛选）
	// 只返回 status=published 的文章
	GetPublishedList(ctx context.Context, req *dto.ArticleListRequest) ([]dto.ArticleListItem, int64, error)

	// GetBySlug 通过 slug 获取文章详情
	// 只返回 status=published 的文章
	// clientIP 用于浏览量去重
	GetBySlug(ctx context.Context, slug string, clientIP string) (*dto.ArticleDetail, error)
}
