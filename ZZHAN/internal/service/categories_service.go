package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type categoriesService struct {
	categoriesRepo repository.CategoriesRepository
}

// NewCategoriesService 创建分类业务实例
func NewCategoriesService(categoriesRepo repository.CategoriesRepository) CategoriesService {
	return &categoriesService{categoriesRepo: categoriesRepo}
}

// GetAll 获取所有分类
func (s *categoriesService) GetAll(ctx context.Context) ([]dto.CategoryItem, error) {
	return s.categoriesRepo.GetAll(ctx)
}
