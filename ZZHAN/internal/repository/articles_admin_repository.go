package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"gorm.io/gorm"
)

type articlesAdminRepository struct {
	db *gorm.DB
}

// NewArticlesAdminRepository 创建后台文章仓储实例
func NewArticlesAdminRepository(db *gorm.DB) ArticlesAdminRepository {
	return &articlesAdminRepository{db: db}
}

// AdminList 获取文章列表
func (r *articlesAdminRepository) AdminList(ctx context.Context, req *dto.AdminArticleListRequest) ([]dto.AdminArticleItem, int64, error) {
	var articles []entity.Article
	var total int64

	query := r.db.WithContext(ctx).Model(&entity.Article{})

	// 关键词搜索（标题）
	if req.Keyword != "" {
		query = query.Where("title LIKE ?", "%"+req.Keyword+"%")
	}

	// 按分类名称筛选
	if req.Category != "" && req.Category != "all" {
		var category entity.Category
		if err := r.db.WithContext(ctx).Where("name = ?", req.Category).First(&category).Error; err == nil {
			query = query.Where("category_id = ?", category.ID)
		}
	}

	// 按标签名称筛选
	if req.Tag != "" && req.Tag != "all" {
		query = query.Where("id IN (?)",
			r.db.Model(&entity.ArticleTag{}).
				Select("article_id").
				Where("tag_id = (?)",
					r.db.Model(&entity.Tag{}).Select("id").Where("name = ?", req.Tag),
				),
		)
	}

	// 按状态筛选
	if req.Status != "" {
		query = query.Where("status = ?", req.Status)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页参数默认值
	page := req.Page
	if page <= 0 {
		page = 1
	}
	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.
		Order("updated_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	result := make([]dto.AdminArticleItem, 0, len(articles))
	for _, article := range articles {
		item := dto.AdminArticleItem{
			ID:         article.ID,
			Slug:       article.Slug,
			Title:      article.Title,
			Summary:    article.Summary,
			CoverImage: article.CoverImage,
			Status:     article.Status,
			Views:      article.Views,
		}

		// 格式化时间
		if article.PublishedAt != nil {
			item.Date = article.PublishedAt.Format("2006-01-02")
		}
		item.Updated = article.UpdatedAt.Format("2006-01-02")

		// 查询分类名称
		var category entity.Category
		if err := r.db.WithContext(ctx).Select("name").First(&category, article.CategoryID).Error; err == nil {
			item.Category = category.Name
		}

		// 查询文章标签
		item.Tags = r.getArticleTagNames(ctx, article.ID)

		result = append(result, item)
	}

	return result, total, nil
}

// AdminGetByID 通过 ID 获取文章详情
func (r *articlesAdminRepository) AdminGetByID(ctx context.Context, id int64) (*dto.AdminArticleDetail, error) {
	var article entity.Article

	if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
		return nil, err
	}

	detail := &dto.AdminArticleDetail{
		ID:         article.ID,
		Slug:       article.Slug,
		Title:      article.Title,
		Summary:    article.Summary,
		CoverImage: article.CoverImage,
		Status:     article.Status,
		Views:      article.Views,
		Content:    article.Content,
	}

	if article.PublishedAt != nil {
		detail.Date = article.PublishedAt.Format("2006-01-02")
	}
	detail.Updated = article.UpdatedAt.Format("2006-01-02")

	// 查询分类名称
	var category entity.Category
	if err := r.db.WithContext(ctx).Select("name").First(&category, article.CategoryID).Error; err == nil {
		detail.Category = category.Name
	}

	// 查询文章标签
	detail.Tags = r.getArticleTagNames(ctx, article.ID)

	return detail, nil
}

// AdminCreate 创建文章
func (r *articlesAdminRepository) AdminCreate(ctx context.Context, req *dto.AdminArticleCreateRequest, authorID int64) (int64, error) {
	// 查询分类 ID
	categoryID, err := r.getCategoryIDByName(ctx, req.Category)
	if err != nil {
		return 0, fmt.Errorf("分类不存在: %s", req.Category)
	}

	// 生成 slug
	slug, err := r.generateSlug(ctx, req.Title)
	if err != nil {
		return 0, err
	}

	// 解析发布时间
	var publishedAt *time.Time
	if req.PublishedAt != "" {
		t, err := time.Parse("2006-01-02", req.PublishedAt)
		if err != nil {
			return 0, fmt.Errorf("发布时间格式错误，应为 2006-01-02")
		}
		publishedAt = &t
	} else {
		now := time.Now()
		publishedAt = &now
	}

	article := entity.Article{
		Title:       req.Title,
		Slug:        slug,
		Summary:     req.Summary,
		CoverImage:  req.CoverImage,
		CategoryID:  categoryID,
		AuthorID:    authorID,
		Content:     req.Content,
		Status:      req.Status,
		PublishedAt: publishedAt,
	}

	// 创建文章
	if err := r.db.WithContext(ctx).Create(&article).Error; err != nil {
		return 0, err
	}

	// 处理标签关联
	if len(req.Tags) > 0 {
		if err := r.syncArticleTags(ctx, article.ID, req.Tags); err != nil {
			return 0, err
		}
	}

	return article.ID, nil
}

// AdminUpdate 更新文章
func (r *articlesAdminRepository) AdminUpdate(ctx context.Context, id int64, req *dto.AdminArticleUpdateRequest) error {
	var article entity.Article
	if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
		return err
	}

	// 查询分类 ID
	categoryID, err := r.getCategoryIDByName(ctx, req.Category)
	if err != nil {
		return fmt.Errorf("分类不存在: %s", req.Category)
	}

	// 解析发布时间
	var publishedAt *time.Time
	if req.PublishedAt != "" {
		t, err := time.Parse("2006-01-02", req.PublishedAt)
		if err != nil {
			return fmt.Errorf("发布时间格式错误，应为 2006-01-02")
		}
		publishedAt = &t
	}

	// 如果标题变了，重新生成 slug
	slug := article.Slug
	if req.Title != article.Title {
		slug, err = r.generateSlug(ctx, req.Title)
		if err != nil {
			return err
		}
	}

	// 更新文章字段
	updates := map[string]interface{}{
		"title":        req.Title,
		"slug":         slug,
		"summary":      req.Summary,
		"cover_image":  req.CoverImage,
		"category_id":  categoryID,
		"content":      req.Content,
		"status":       req.Status,
		"published_at": publishedAt,
	}

	if err := r.db.WithContext(ctx).Model(&article).Updates(updates).Error; err != nil {
		return err
	}

	// 同步标签关联
	if err := r.syncArticleTags(ctx, id, req.Tags); err != nil {
		return err
	}

	return nil
}

// AdminDelete 删除文章（软删除）
func (r *articlesAdminRepository) AdminDelete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&entity.Article{}, id).Error
}

// AdminUpdateStatus 修改文章状态
func (r *articlesAdminRepository) AdminUpdateStatus(ctx context.Context, id int64, status string) error {
	// 发布前校验必填信息
	if status == "published" {
		var article entity.Article
		if err := r.db.WithContext(ctx).First(&article, id).Error; err != nil {
			return fmt.Errorf("文章不存在")
		}

		missing := []string{}
		if article.CategoryID == 0 {
			missing = append(missing, "分类")
		}
		if article.CoverImage == "" {
			missing = append(missing, "封面")
		}
		// 检查标签
		var tagCount int64
		r.db.WithContext(ctx).Model(&entity.ArticleTag{}).Where("article_id = ?", id).Count(&tagCount)
		if tagCount == 0 {
			missing = append(missing, "标签")
		}
		if len(missing) > 0 {
			return fmt.Errorf("文章缺少%s，请先编辑补充", strings.Join(missing, "、"))
		}
	}

	updates := map[string]interface{}{
		"status": status,
	}

	// 如果是发布状态且之前没有发布时间，设置发布时间
	if status == "published" {
		var article entity.Article
		if err := r.db.WithContext(ctx).Select("published_at").First(&article, id).Error; err == nil {
			if article.PublishedAt == nil {
				now := time.Now()
				updates["published_at"] = now
			}
		}
	}

	return r.db.WithContext(ctx).Model(&entity.Article{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// getCategoryIDByName 通过分类名称获取分类 ID
func (r *articlesAdminRepository) getCategoryIDByName(ctx context.Context, name string) (int64, error) {
	var category entity.Category
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&category).Error; err != nil {
		return 0, err
	}
	return category.ID, nil
}

// getArticleTagNames 获取文章的标签名称列表
func (r *articlesAdminRepository) getArticleTagNames(ctx context.Context, articleID int64) []string {
	var tags []string

	r.db.WithContext(ctx).
		Table("tags").
		Select("tags.name").
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Pluck("name", &tags)

	if tags == nil {
		tags = []string{}
	}

	return tags
}

// syncArticleTags 同步文章标签关联
func (r *articlesAdminRepository) syncArticleTags(ctx context.Context, articleID int64, tagNames []string) error {
	// 删除现有关联
	if err := r.db.WithContext(ctx).Where("article_id = ?", articleID).Delete(&entity.ArticleTag{}).Error; err != nil {
		return err
	}

	if len(tagNames) == 0 {
		return nil
	}

	// 查找或创建标签，建立关联
	for _, name := range tagNames {
		name = strings.TrimSpace(name)
		if name == "" {
			continue
		}

		var tag entity.Tag
		// 查找已有标签
		if err := r.db.WithContext(ctx).Where("name = ?", name).First(&tag).Error; err != nil {
			// 标签不存在，创建新标签
			tag = entity.Tag{
				Name: name,
				Slug: generateTagSlug(name),
			}
			if err := r.db.WithContext(ctx).Create(&tag).Error; err != nil {
				return err
			}
		}

		// 创建关联
		articleTag := entity.ArticleTag{
			ArticleID: articleID,
			TagID:     tag.ID,
		}
		if err := r.db.WithContext(ctx).Create(&articleTag).Error; err != nil {
			return err
		}
	}

	return nil
}

// generateSlug 生成 URL 友好的 slug，冲突时追加后缀
func (r *articlesAdminRepository) generateSlug(ctx context.Context, title string) (string, error) {
	slug := slugify(title)
	if slug == "" {
		slug = "article"
	}

	baseSlug := slug
	for i := 1; ; i++ {
		var count int64
		if err := r.db.WithContext(ctx).Model(&entity.Article{}).
			Where("slug = ?", slug).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return slug, nil
		}
		slug = fmt.Sprintf("%s-%d", baseSlug, i)
	}
}

// slugify 将标题转为 URL 友好格式
func slugify(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	// 保留中文、字母、数字，其余替换为连字符
	reg := regexp.MustCompile(`[^\p{Han}\p{L}\p{N}]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 200 {
		s = s[:200]
	}
	return s
}

// generateTagSlug 为标签生成 slug
func generateTagSlug(name string) string {
	s := strings.ToLower(strings.TrimSpace(name))
	reg := regexp.MustCompile(`[^\p{Han}\p{L}\p{N}]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if s == "" {
		s = "tag"
	}
	return s
}
