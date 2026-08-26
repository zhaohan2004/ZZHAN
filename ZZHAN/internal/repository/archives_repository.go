package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"
	"fmt"
	"sort"

	"gorm.io/gorm"
)

type archivesRepository struct {
	db *gorm.DB
}

// NewArchivesRepository 创建归档仓储实例
func NewArchivesRepository(db *gorm.DB) ArchivesRepository {
	return &archivesRepository{db: db}
}

// GetArchives 获取归档列表（按年月分组）
func (r *archivesRepository) GetArchives(ctx context.Context) ([]dto.ArchiveItem, error) {
	var articles []entity.Article

	// 查询已发布文章，按发布时间倒序
	if err := r.db.WithContext(ctx).
		Where("status = ?", "published").
		Order("published_at DESC").
		Find(&articles).Error; err != nil {
		return nil, err
	}

	// 预加载分类名称
	categoryMap := make(map[int64]string)
	var categories []entity.Category
	if err := r.db.WithContext(ctx).Select("id, name").Find(&categories).Error; err == nil {
		for _, c := range categories {
			categoryMap[c.ID] = c.Name
		}
	}

	// 按年月分组
	groupMap := make(map[string]*dto.ArchiveItem)
	var keys []string

	for _, article := range articles {
		if article.PublishedAt == nil {
			continue
		}
		year := article.PublishedAt.Year()
		month := int(article.PublishedAt.Month())
		key := fmt.Sprintf("%04d-%02d", year, month)

		if _, ok := groupMap[key]; !ok {
			groupMap[key] = &dto.ArchiveItem{
				Year:     fmt.Sprintf("%d", year),
				Month:    fmt.Sprintf("%d", month),
				Articles: []dto.ArchiveArticle{},
			}
			keys = append(keys, key)
		}

		item := groupMap[key]
		item.Count++
		item.Articles = append(item.Articles, dto.ArchiveArticle{
			ID:       article.ID,
			Slug:     article.Slug,
			Title:    article.Title,
			Date:     article.PublishedAt.Format("2006-01-02"),
			Category: categoryMap[article.CategoryID],
			Views:    article.Views,
		})
	}

	// 按年月倒序排列
	sort.Sort(sort.Reverse(sort.StringSlice(keys)))

	result := make([]dto.ArchiveItem, 0, len(keys))
	for _, key := range keys {
		result = append(result, *groupMap[key])
	}

	return result, nil
}
