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

	// GetBySlug 通过 slug 获取文章详情  clientIP 用于浏览量去重  userID 用于查询点赞状态（未登录传 0）
	GetBySlug(ctx context.Context, slug string, clientIP string, userID int64) (*dto.ArticleDetail, error)
}
