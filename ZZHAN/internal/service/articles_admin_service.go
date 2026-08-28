package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type articlesAdminService struct {
	articlesAdminRepo repository.ArticlesAdminRepository
}

// NewArticlesAdminService 创建后台文章业务实例
func NewArticlesAdminService(articlesAdminRepo repository.ArticlesAdminRepository) ArticlesAdminService {
	return &articlesAdminService{articlesAdminRepo: articlesAdminRepo}
}

// AdminList 获取文章列表
func (s *articlesAdminService) AdminList(ctx context.Context, req *dto.AdminArticleListRequest) (*dto.AdminArticleListResponse, error) {
	// 分页默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	articles, total, err := s.articlesAdminRepo.AdminList(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.AdminArticleListResponse{
		List:  articles,
		Total: total,
	}, nil
}

// AdminGetByID 获取文章详情
func (s *articlesAdminService) AdminGetByID(ctx context.Context, id int64) (*dto.AdminArticleDetail, error) {
	return s.articlesAdminRepo.AdminGetByID(ctx, id)
}

// AdminCreate 创建文章
func (s *articlesAdminService) AdminCreate(ctx context.Context, req *dto.AdminArticleCreateRequest, authorID int64) (*dto.AdminArticleDetail, error) {
	id, err := s.articlesAdminRepo.AdminCreate(ctx, req, authorID)
	if err != nil {
		return nil, err
	}

	return s.articlesAdminRepo.AdminGetByID(ctx, id)
}

// AdminUpdate 更新文章
func (s *articlesAdminService) AdminUpdate(ctx context.Context, id int64, req *dto.AdminArticleUpdateRequest) (*dto.AdminArticleDetail, error) {
	if err := s.articlesAdminRepo.AdminUpdate(ctx, id, req); err != nil {
		return nil, err
	}

	return s.articlesAdminRepo.AdminGetByID(ctx, id)
}

// AdminDelete 删除文章
func (s *articlesAdminService) AdminDelete(ctx context.Context, id int64) error {
	return s.articlesAdminRepo.AdminDelete(ctx, id)
}

// AdminUpdateStatus 修改文章状态
func (s *articlesAdminService) AdminUpdateStatus(ctx context.Context, id int64, status string) (*dto.AdminArticleDetail, error) {
	if err := s.articlesAdminRepo.AdminUpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return s.articlesAdminRepo.AdminGetByID(ctx, id)
}
