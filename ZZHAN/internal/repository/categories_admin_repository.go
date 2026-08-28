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

type categoriesAdminRepository struct {
	db *gorm.DB
}

// NewCategoriesAdminRepository 创建后台分类仓储实例
func NewCategoriesAdminRepository(db *gorm.DB) CategoriesAdminRepository {
	return &categoriesAdminRepository{db: db}
}

// AdminList 分页查询分类列表
func (r *categoriesAdminRepository) AdminList(ctx context.Context, req *dto.AdminCategoryListRequest) ([]dto.AdminCategoryItem, int64, error) {
	var total int64

	// 文章数统计子查询
	countSub := r.db.Model(&entity.Article{}).
		Select("COUNT(*)").
		Where("articles.category_id = categories.id AND articles.status = ? AND articles.deleted_at IS NULL", "published")

	// 基础查询
	base := r.db.WithContext(ctx).Table("categories")

	if req.Keyword != "" {
		base = base.Where("categories.name LIKE ?", "%"+req.Keyword+"%")
	}
	if req.Status != "" && req.Status != "all" {
		base = base.Where("categories.status = ?", req.Status)
	}

	// 需要按文章数过滤时，用 JOIN 子查询
	if req.MinCount != nil || req.MaxCount != nil {
		// 构建分类+文章数的派生表
		derived := r.db.Table("categories").
			Select("categories.id AS cid, COALESCE(ac.cnt, 0) AS article_count").
			Joins("LEFT JOIN (?) ac ON ac.category_id = categories.id",
				r.db.Model(&entity.Article{}).
					Select("category_id, COUNT(*) AS cnt").
					Where("status = ? AND deleted_at IS NULL", "published").
					Group("category_id"),
			)

		if req.Keyword != "" {
			derived = derived.Where("categories.name LIKE ?", "%"+req.Keyword+"%")
		}
		if req.Status != "" && req.Status != "all" {
			derived = derived.Where("categories.status = ?", req.Status)
		}
		if req.MinCount != nil {
			derived = derived.Having("article_count >= ?", *req.MinCount)
		}
		if req.MaxCount != nil {
			derived = derived.Having("article_count <= ?", *req.MaxCount)
		}

		// 统计总数
		if err := r.db.WithContext(ctx).Table("(?) AS sub", derived).Count(&total).Error; err != nil {
			return nil, 0, err
		}

		// 分页参数
		page := req.Page
		if page <= 0 {
			page = 1
		}
		pageSize := req.PageSize
		if pageSize <= 0 {
			pageSize = 10
		}
		offset := (page - 1) * pageSize

		// 查询 ID 列表
		var ids []int64
		if err := derived.Select("cid").Order("cid ASC").Offset(offset).Limit(pageSize).Scan(&ids).Error; err != nil {
			return nil, 0, err
		}

		if len(ids) == 0 {
			return []dto.AdminCategoryItem{}, 0, nil
		}

		// 查询分类详情
		var categories []entity.Category
		if err := r.db.WithContext(ctx).Where("id IN ?", ids).Find(&categories).Error; err != nil {
			return nil, 0, err
		}

		// 按 ID 顺序排列
		idOrder := make(map[int64]int, len(ids))
		for i, id := range ids {
			idOrder[id] = i
		}
		ordered := make([]entity.Category, len(categories))
		for _, c := range categories {
			ordered[idOrder[c.ID]] = c
		}

		// 转换为 DTO
		result := make([]dto.AdminCategoryItem, 0, len(ordered))
		for _, c := range ordered {
			var articleCount int64
			r.db.WithContext(ctx).Model(&entity.Article{}).
				Where("category_id = ? AND status = ? AND deleted_at IS NULL", c.ID, "published").
				Count(&articleCount)

			result = append(result, dto.AdminCategoryItem{
				ID:           c.ID,
				Name:         c.Name,
				Slug:         c.Slug,
				Icon:         c.Icon,
				Description:  c.Description,
				Color:        c.Color,
				Status:       c.Status,
				ArticleCount: articleCount,
				CreatedAt:    c.CreatedAt.Format("2006-01-02 15:04"),
				UpdatedAt:    c.UpdatedAt.Format("2006-01-02 15:04"),
			})
		}

		return result, total, nil
	}

	// 无文章数过滤时的简单查询
	if err := base.Count(&total).Error; err != nil {
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

	var categories []entity.Category
	if err := r.db.WithContext(ctx).
		Select("categories.*, (?) AS article_count", countSub).
		Order("id ASC").
		Offset(offset).
		Limit(pageSize).
		Find(&categories).Error; err != nil {
		return nil, 0, err
	}

	result := make([]dto.AdminCategoryItem, 0, len(categories))
	for _, c := range categories {
		var articleCount int64
		r.db.WithContext(ctx).Model(&entity.Article{}).
			Where("category_id = ? AND status = ? AND deleted_at IS NULL", c.ID, "published").
			Count(&articleCount)

		result = append(result, dto.AdminCategoryItem{
			ID:           c.ID,
			Name:         c.Name,
			Slug:         c.Slug,
			Icon:         c.Icon,
			Description:  c.Description,
			Color:        c.Color,
			Status:       c.Status,
			ArticleCount: articleCount,
			CreatedAt:    c.CreatedAt.Format("2006-01-02 15:04"),
			UpdatedAt:    c.UpdatedAt.Format("2006-01-02 15:04"),
		})
	}

	return result, total, nil
}

// AdminGetByID 通过 ID 获取分类
func (r *categoriesAdminRepository) AdminGetByID(ctx context.Context, id int64) (*dto.AdminCategoryItem, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).First(&category, id).Error; err != nil {
		return nil, err
	}

	// 统计文章数
	var count int64
	r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("category_id = ? AND status = ? AND deleted_at IS NULL", id, "published").
		Count(&count)

	return &dto.AdminCategoryItem{
		ID:           category.ID,
		Name:         category.Name,
		Slug:         category.Slug,
		Icon:         category.Icon,
		Description:  category.Description,
		Color:        category.Color,
		Status:       category.Status,
		ArticleCount: count,
		CreatedAt:    category.CreatedAt.Format("2006-01-02 15:04"),
		UpdatedAt:    category.UpdatedAt.Format("2006-01-02 15:04"),
	}, nil
}

// AdminCreate 创建分类
func (r *categoriesAdminRepository) AdminCreate(ctx context.Context, name, slug, desc, color, status string) (int64, error) {
	// 自动生成 slug
	if slug == "" {
		slug = r.generateSlug(ctx, name)
	}

	if color == "" {
		color = "#3b82f6"
	}
	if status == "" {
		status = "active"
	}

	category := entity.Category{
		Name:        name,
		Slug:        slug,
		Description: desc,
		Color:       color,
		Status:      status,
	}

	if err := r.db.WithContext(ctx).Create(&category).Error; err != nil {
		return 0, err
	}

	return category.ID, nil
}

// AdminUpdate 更新分类
func (r *categoriesAdminRepository) AdminUpdate(ctx context.Context, id int64, updates map[string]interface{}) error {
	return r.db.WithContext(ctx).Model(&entity.Category{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// AdminUpdateStatus 修改分类状态
func (r *categoriesAdminRepository) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	return r.db.WithContext(ctx).Model(&entity.Category{}).
		Where("id = ?", id).
		Update("status", status).Error
}

// AdminDelete 删除分类
func (r *categoriesAdminRepository) AdminDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Category{}, id).Error
}

// generateSlug 生成 URL 友好的 slug，冲突时追加后缀
func (r *categoriesAdminRepository) generateSlug(ctx context.Context, name string) string {
	slug := categorySlugify(name)
	if slug == "" {
		slug = "category"
	}

	baseSlug := slug
	for i := 1; ; i++ {
		var count int64
		r.db.WithContext(ctx).Model(&entity.Category{}).
			Where("slug = ?", slug).Count(&count)
		if count == 0 {
			return slug
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}
}

// categorySlugify 将名称转为 URL 友好格式
func categorySlugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	reg := regexp.MustCompile(`[^\p{Han}\p{L}\p{N}]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}
