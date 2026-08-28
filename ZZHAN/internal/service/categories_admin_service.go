package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type categoriesAdminService struct {
	categoriesAdminRepo repository.CategoriesAdminRepository
}

// NewCategoriesAdminService 创建后台分类业务实例
func NewCategoriesAdminService(categoriesAdminRepo repository.CategoriesAdminRepository) CategoriesAdminService {
	return &categoriesAdminService{categoriesAdminRepo: categoriesAdminRepo}
}

// AdminList 获取分类列表
func (s *categoriesAdminService) AdminList(ctx context.Context, req *dto.AdminCategoryListRequest) (*dto.AdminCategoryListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := s.categoriesAdminRepo.AdminList(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.AdminCategoryListResponse{
		List:  list,
		Total: total,
	}, nil
}

// AdminGetByID 获取分类详情
func (s *categoriesAdminService) AdminGetByID(ctx context.Context, id int64) (*dto.AdminCategoryItem, error) {
	return s.categoriesAdminRepo.AdminGetByID(ctx, id)
}

// AdminCreate 创建分类
func (s *categoriesAdminService) AdminCreate(ctx context.Context, req *dto.AdminCategoryCreateRequest) (*dto.AdminCategoryItem, error) {
	id, err := s.categoriesAdminRepo.AdminCreate(ctx, req.Name, req.Slug, req.Desc, req.Color, req.Status)
	if err != nil {
		return nil, err
	}

	return s.categoriesAdminRepo.AdminGetByID(ctx, id)
}

// AdminUpdate 更新分类
func (s *categoriesAdminService) AdminUpdate(ctx context.Context, id int64, req *dto.AdminCategoryUpdateRequest) (*dto.AdminCategoryItem, error) {
	updates := make(map[string]interface{})
	if req.Name != nil {
		updates["name"] = *req.Name
	}
	if req.Slug != nil {
		updates["slug"] = *req.Slug
	}
	if req.Desc != nil {
		updates["description"] = *req.Desc
	}
	if req.Color != nil {
		updates["color"] = *req.Color
	}
	if req.Status != nil {
		updates["status"] = *req.Status
	}

	if len(updates) == 0 {
		return s.categoriesAdminRepo.AdminGetByID(ctx, id)
	}

	if err := s.categoriesAdminRepo.AdminUpdate(ctx, id, updates); err != nil {
		return nil, err
	}

	return s.categoriesAdminRepo.AdminGetByID(ctx, id)
}

// AdminUpdateStatus 修改分类状态
func (s *categoriesAdminService) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.categoriesAdminRepo.AdminUpdateStatus(ctx, id, status)
}

// AdminDelete 删除分类
func (s *categoriesAdminService) AdminDelete(ctx context.Context, id int64) error {
	return s.categoriesAdminRepo.AdminDelete(ctx, id)
}
