package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CategoriesRepository 分类仓储接口
type CategoriesRepository interface {
	// GetAll 获取所有分类（含文章数量）
	GetAll(ctx context.Context) ([]dto.CategoryItem, error)
}
