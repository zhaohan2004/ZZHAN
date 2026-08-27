package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/pkg/response"
	"context"
)

type articlesService struct {
	articlesRepo repository.ArticlesRepository
}

// NewArticlesService 创建文章业务实例
func NewArticlesService(articlesRepo repository.ArticlesRepository) ArticlesService {
	return &articlesService{articlesRepo: articlesRepo}
}

// GetPublishedList 获取已发布文章列表
func (s *articlesService) GetPublishedList(ctx context.Context, req *dto.ArticleListRequest) (*response.PageResponse, error) {
	// 调用 repository 获取数据
	articles, total, err := s.articlesRepo.GetPublishedList(ctx, req)
	if err != nil {
		return nil, err
	}

	// 使用 pkg/response 的 NewPageResponse 构建分页响应
	return response.NewPageResponse(articles, total, req.Page, req.Size), nil
}

// GetBySlug 通过 slug 获取文章详情
func (s *articlesService) GetBySlug(ctx context.Context, slug string, clientIP string) (*dto.ArticleDetail, error) {
	return s.articlesRepo.GetBySlug(ctx, slug, clientIP)
}
