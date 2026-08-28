package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"
	"fmt"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

type tagsAdminRepository struct {
	db *gorm.DB
}

// NewTagsAdminRepository 创建后台标签仓储实例
func NewTagsAdminRepository(db *gorm.DB) TagsAdminRepository {
	return &tagsAdminRepository{db: db}
}

// AdminList 分页查询标签列表
func (r *tagsAdminRepository) AdminList(ctx context.Context, req *dto.AdminTagListRequest) ([]dto.AdminTagItem, int64, error) {
	var total int64

	// 文章数统计子查询（通过 article_tags 关联已发布文章）
	countSub := r.db.Table("article_tags").
		Select("COUNT(*)").
		Joins("JOIN articles ON articles.id = article_tags.article_id AND articles.status = ? AND articles.deleted_at IS NULL", "published").
		Where("article_tags.tag_id = tags.id")

	// 需要按文章数过滤时
	if req.MinCount != nil || req.MaxCount != nil {
		derived := r.db.Table("tags").
			Select("tags.id AS tid, COALESCE(ac.cnt, 0) AS article_count").
			Joins("LEFT JOIN (?) ac ON ac.tag_id = tags.id",
				r.db.Table("article_tags").
					Select("tag_id, COUNT(*) AS cnt").
					Joins("JOIN articles ON articles.id = article_tags.article_id AND articles.status = ? AND articles.deleted_at IS NULL", "published").
					Group("tag_id"),
			)

		if req.Keyword != "" {
			derived = derived.Where("tags.name LIKE ?", "%"+req.Keyword+"%")
		}
		if req.Status != "" && req.Status != "all" {
			derived = derived.Where("tags.status = ?", req.Status)
		}
		if req.MinCount != nil {
			derived = derived.Having("article_count >= ?", *req.MinCount)
		}
		if req.MaxCount != nil {
			derived = derived.Having("article_count <= ?", *req.MaxCount)
		}

		if err := r.db.WithContext(ctx).Table("(?) AS sub", derived).Count(&total).Error; err != nil {
			return nil, 0, err
		}

		page := req.Page
		if page <= 0 {
			page = 1
		}
		pageSize := req.PageSize
		if pageSize <= 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize

		var ids []int64
		if err := derived.Select("tid").Order("tid ASC").Offset(offset).Limit(pageSize).Scan(&ids).Error; err != nil {
			return nil, 0, err
		}

		if len(ids) == 0 {
			return []dto.AdminTagItem{}, 0, nil
		}

		var tags []entity.Tag
		if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&tags).Error; err != nil {
			return nil, 0, err
		}

		idOrder := make(map[int64]int, len(ids))
		for i, id := range ids {
			idOrder[id] = i
		}
		ordered := make([]entity.Tag, len(tags))
		for _, t := range tags {
			ordered[idOrder[t.ID]] = t
		}

		result := make([]dto.AdminTagItem, 0, len(ordered))
		for _, t := range ordered {
			var articleCount int64
			r.db.WithContext(ctx).Table("article_tags").
				Joins("JOIN articles ON articles.id = article_tags.article_id AND articles.status = ? AND articles.deleted_at IS NULL", "published").
				Where("article_tags.tag_id = ?", t.ID).
				Count(&articleCount)

			result = append(result, dto.AdminTagItem{
				ID:           t.ID,
				Name:         t.Name,
				Status:       t.Status,
				ArticleCount: articleCount,
				CreatedAt:    t.CreatedAt.Format("2006-01-02 15:04"),
				UpdatedAt:    t.UpdatedAt.Format("2006-01-02 15:04"),
			})
		}

		return result, total, nil
	}

	// 无文章数过滤时
	query := r.db.WithContext(ctx).Model(&entity.Tag{})
	if req.Keyword != "" {
		query = query.Where("name LIKE ?", "%"+req.Keyword+"%")
	}
	if req.Status != "" && req.Status != "all" {
		query = query.Where("status = ?", req.Status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	offset := (page - 1) * pageSize

	var tags []entity.Tag
	if err := query.
		Select("tags.*, (?) AS article_count", countSub).
		Order("id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&tags).Error; err != nil {
		return nil, 0, err
	}

	result := make([]dto.AdminTagItem, 0, len(tags))
	for _, t := range tags {
		var articleCount int64
		r.db.WithContext(ctx).Table("article_tags").
			Joins("JOIN articles ON articles.id = article_tags.article_id AND articles.status = ? AND articles.deleted_at IS NULL", "published").
			Where("article_tags.tag_id = ?", t.ID).
			Count(&articleCount)

		result = append(result, dto.AdminTagItem{
			ID:           t.ID,
			Name:         t.Name,
			Status:       t.Status,
			ArticleCount: articleCount,
			CreatedAt:    t.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt:    t.UpdatedAt.Format("2006-01-02 15:04"),
		})
	}

	return result, total, nil
}

// AdminGetByID 通过 ID 获取标签
func (r *tagsAdminRepository) AdminGetByID(ctx context.Context, id int64) (*dto.AdminTagItem, error) {
	var tag entity.Tag
	if err := r.db.WithContext(ctx).First(&tag, id).Error; err != nil {
		return nil, err
	}

	var articleCount int64
	r.db.WithContext(ctx).Table("article_tags").
		Joins("JOIN articles ON articles.id = article_tags.article_id AND articles.status = ? AND articles.deleted_at IS NULL", "published").
		Where("article_tags.tag_id = ?", id).
		Count(&articleCount)

	return &dto.AdminTagItem{
		ID:           tag.ID,
		Name:         tag.Name,
		Status:       tag.Status,
		ArticleCount: articleCount,
		CreatedAt:    tag.CreatedAt.Format("2006-01-02 15:04"),
		UpdatedAt:    tag.UpdatedAt.Format("2006-01-02 15:04"),
	}, nil
}

// AdminCreate 创建标签
func (r *tagsAdminRepository) AdminCreate(ctx context.Context, name, status string) (int64, error) {
	slug := tagSlugify(name)
	if slug == "" {
		slug = "tag"
	}

	// 检查 slug 唯一性
	baseSlug := slug
	for i := 1; ; i++ {
		var count int64
		r.db.WithContext(ctx).Model(&entity.Tag{}).Where("slug = ?", slug).Count(&count)
		if count == 0 {
			break
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}

	if status == "" {
		status = "active"
	}

	tag := entity.Tag{
		Name:   name,
		Slug:   slug,
		Status: status,
	}

	if err := r.db.WithContext(ctx).Create(&tag).Error; err != nil {
		return 0, err
	}

	return tag.ID, nil
}

// AdminUpdate 更新标签
func (r *tagsAdminRepository) AdminUpdate(ctx context.Context, id int64, name string) error {
	return r.db.WithContext(ctx).Model(&entity.Tag{}).
		Where("id = ?", id).
		Update("name", name).Error
}

// AdminUpdateStatus 修改标签状态
func (r *tagsAdminRepository) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Tag{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// AdminDelete 删除标签（同时清理 article_tags 关联）
func (r *tagsAdminRepository) AdminDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 删除关联
		if err := tx.Where("tag_id = ?", id).Delete(&entity.ArticleTag{}).Error; err != nil {
			return err
		}
		// 删除标签
		return tx.Delete(&entity.Tag{}, id).Error
	})
}

// tagSlugify 将标签名转为 URL 友好格式
func tagSlugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	reg := regexp.MustCompile(`[^\p{Han}\p{L}\p{N}]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}
