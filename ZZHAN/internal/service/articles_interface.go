package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/pkg/response"
	"context"
)

// ArticlesService 文章业务接口
type ArticlesService interface {
	// GetPublishedList 获取已发布文章列表
	GetPublishedList(ctx context.Context, req *dto.ArticleListRequest) (*response.PageResponse, error)

	// GetBySlug 通过 slug 获取文章详情
	GetBySlug(ctx context.Context, slug string) (*dto.ArticleDetail, error)
}
