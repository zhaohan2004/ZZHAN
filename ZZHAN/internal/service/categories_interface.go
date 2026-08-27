package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CategoriesService 分类业务接口
type CategoriesService interface {
	// GetAll 获取所有分类
	GetAll(ctx context.Context) ([]dto.CategoryItem, error)
}
