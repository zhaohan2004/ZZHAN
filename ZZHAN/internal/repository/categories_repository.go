package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"

	"gorm.io/gorm"
)

type categoriesRepository struct {
	db *gorm.DB
}

// NewCategoriesRepository 创建分类仓储实例
func NewCategoriesRepository(db *gorm.DB) CategoriesRepository {
	return &categoriesRepository{db: db}
}

// GetAll 获取所有分类（含文章数量）
func (r *categoriesRepository) GetAll(ctx context.Context) ([]dto.CategoryItem, error) {
	var categories []entity.Category

	// 查询所有分类，按排序值升序
	if err := r.db.WithContext(ctx).
		Order("sort_order ASC, id ASC").
		Find(&categories).Error; err != nil {
		return nil, err
	}

	// 转换为 DTO
	result := make([]dto.CategoryItem, 0, len(categories))
	for _, cat := range categories {
		item := dto.CategoryItem{
			ID:          cat.ID,
			Name:        cat.Name,
			Slug:        cat.Slug,
			Icon:        cat.Icon,
			Description: cat.Description,
			Color:       cat.Color,
		}

		// 统计该分类下的已发布文章数量
		var count int64
		r.db.WithContext(ctx).
			Model(&entity.Article{}).
			Where("category_id = ? AND status = ?", cat.ID, "published").
			Count(&count)
		item.ArticleCount = count

		result = append(result, item)
	}

	return result, nil
}
