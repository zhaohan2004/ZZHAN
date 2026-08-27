package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type tagsService struct {
	tagsRepo repository.TagsRepository
}

// NewTagsService 创建标签业务实例
func NewTagsService(tagsRepo repository.TagsRepository) TagsService {
	return &tagsService{tagsRepo: tagsRepo}
}

// GetAll 获取所有标签
func (s *tagsService) GetAll(ctx context.Context) ([]dto.TagListItem, error) {
	return s.tagsRepo.GetAll(ctx)
}
