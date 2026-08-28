package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"

	"gorm.io/gorm"
)

type tagsRepository struct {
	db *gorm.DB
}

// NewTagsRepository 创建标签仓储实例
func NewTagsRepository(db *gorm.DB) TagsRepository {
	return &tagsRepository{db: db}
}

// GetAll 获取所有标签（含文章数量）
func (r *tagsRepository) GetAll(ctx context.Context) ([]dto.TagListItem, error) {
	var tags []entity.Tag

	// 查询所有启用的标签
	if err := r.db.WithContext(ctx).
		Where("status = ?", "active").
		Order("id ASC").
		Find(&tags).Error; err != nil {
		return nil, err
	}

	// 转换为 DTO
	result := make([]dto.TagListItem, 0, len(tags))
	for _, tag := range tags {
		item := dto.TagListItem{
			ID:   tag.ID,
			Name: tag.Name,
		}

		// 统计该标签下的已发布文章数量（通过关联表）
		var count int64
		r.db.WithContext(ctx).
			Model(&entity.ArticleTag{}).
			Joins("JOIN articles ON articles.id = article_tags.article_id").
			Where("article_tags.tag_id = ? AND articles.status = ?", tag.ID, "published").
			Count(&count)
		item.ArticleCount = count

		result = append(result, item)
	}

	return result, nil
}
