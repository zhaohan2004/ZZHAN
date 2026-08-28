package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type tagsAdminService struct {
	tagsAdminRepo repository.TagsAdminRepository
}

// NewTagsAdminService 创建后台标签业务实例
func NewTagsAdminService(tagsAdminRepo repository.TagsAdminRepository) TagsAdminService {
	return &tagsAdminService{tagsAdminRepo: tagsAdminRepo}
}

// AdminList 获取标签列表
func (s *tagsAdminService) AdminList(ctx context.Context, req *dto.AdminTagListRequest) (*dto.AdminTagListResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	list, total, err := s.tagsAdminRepo.AdminList(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.AdminTagListResponse{
		List:  list,
		Total: total,
	}, nil
}

// AdminGetByID 获取标签详情
func (s *tagsAdminService) AdminGetByID(ctx context.Context, id int64) (*dto.AdminTagItem, error) {
	return s.tagsAdminRepo.AdminGetByID(ctx, id)
}

// AdminCreate 创建标签
func (s *tagsAdminService) AdminCreate(ctx context.Context, req *dto.AdminTagCreateRequest) (*dto.AdminTagItem, error) {
	id, err := s.tagsAdminRepo.AdminCreate(ctx, req.Name, req.Status)
	if err != nil {
		return nil, err
	}

	return s.tagsAdminRepo.AdminGetByID(ctx, id)
}

// AdminUpdate 更新标签
func (s *tagsAdminService) AdminUpdate(ctx context.Context, id int64, req *dto.AdminTagUpdateRequest) (*dto.AdminTagItem, error) {
	if req.Name != nil {
		if err := s.tagsAdminRepo.AdminUpdate(ctx, id, *req.Name); err != nil {
			return nil, err
		}
	}

	return s.tagsAdminRepo.AdminGetByID(ctx, id)
}

// AdminUpdateStatus 修改标签状态
func (s *tagsAdminService) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return s.tagsAdminRepo.AdminUpdateStatus(ctx, id, status)
}

// AdminDelete 删除标签
func (s *tagsAdminService) AdminDelete(ctx context.Context, id int64) error {
	return s.tagsAdminRepo.AdminDelete(ctx, id)
}
