package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"

	"gorm.io/gorm"
)

type articlesRepository struct {
	db *gorm.DB
}

// NewArticlesRepository 创建文章仓储实例
func NewArticlesRepository(db *gorm.DB) ArticlesRepository {
	return &articlesRepository{db: db}
}

// GetPublishedList 获取已发布文章列表
func (r *articlesRepository) GetPublishedList(ctx context.Context, req *dto.ArticleListRequest) ([]dto.ArticleListItem, int64, error) {
	var articles []entity.Article
	var total int64

	// 构建基础查询：只查已发布文章
	query := r.db.WithContext(ctx).
		Model(&entity.Article{}).
		Where("status = ?", "published")

	// 按分类筛选
	if req.CategoryID > 0 {
		query = query.Where("category_id = ?", req.CategoryID)
	}

	// 按标签筛选（需要子查询）
	if req.TagID > 0 {
		query = query.Where("id IN (?)",
			r.db.Model(&entity.ArticleTag{}).
				Select("article_id").
				Where("tag_id = ?", req.TagID),
		)
	}

	// 关键词搜索（标题或摘要）
	if req.Keyword != "" {
		keyword := "%" + req.Keyword + "%"
		query = query.Where("(title LIKE ? OR summary LIKE ?)", keyword, keyword)
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询，按发布时间倒序
	offset := (req.Page - 1) * req.Size
	if err := query.
		Order("published_at DESC").
		Offset(offset).
		Limit(req.Size).
		Find(&articles).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	result := make([]dto.ArticleListItem, 0, len(articles))
	for _, article := range articles {
		item := dto.ArticleListItem{
			ID:           article.ID,
			Title:        article.Title,
			Slug:         article.Slug,
			Summary:      article.Summary,
			CoverImage:   article.CoverImage,
			CategoryID:   article.CategoryID,
			Views:        article.Views,
			Likes:        article.Likes,
			CommentCount: article.CommentCount,
			PublishedAt:  *article.PublishedAt,
		}

		// 查询分类名称
		var category entity.Category
		if err := r.db.WithContext(ctx).
			Select("name").
			First(&category, article.CategoryID).Error; err == nil {
			item.CategoryName = category.Name
		}

		// 查询文章标签
		item.Tags = r.getArticleTags(ctx, article.ID)

		// 查询作者名称
		var admin entity.Admin
		if err := r.db.WithContext(ctx).
			Select("username").
			First(&admin, article.AuthorID).Error; err == nil {
			item.AuthorName = admin.Username
		}

		result = append(result, item)
	}

	return result, total, nil
}

// GetBySlug 通过 slug 获取文章详情
func (r *articlesRepository) GetBySlug(ctx context.Context, slug string) (*dto.ArticleDetail, error) {
	var article entity.Article

	// 查询文章
	if err := r.db.WithContext(ctx).
		Where("slug = ? AND status = ?", slug, "published").
		First(&article).Error; err != nil {
		return nil, err
	}

	// 构建详情 DTO
	detail := &dto.ArticleDetail{
		ID:           article.ID,
		Title:        article.Title,
		Slug:         article.Slug,
		Summary:      article.Summary,
		CoverImage:   article.CoverImage,
		CategoryID:   article.CategoryID,
		Content:      article.Content,
		Views:        article.Views,
		Likes:        article.Likes,
		CommentCount: article.CommentCount,
		PublishedAt:  *article.PublishedAt,
	}

	// 查询分类名称
	var category entity.Category
	if err := r.db.WithContext(ctx).
		Select("name").
		First(&category, article.CategoryID).Error; err == nil {
		detail.CategoryName = category.Name
	}

	// 查询文章标签
	detail.Tags = r.getArticleTags(ctx, article.ID)

	// 查询作者名称
	var admin entity.Admin
	if err := r.db.WithContext(ctx).
		Select("username").
		First(&admin, article.AuthorID).Error; err == nil {
		detail.AuthorName = admin.Username
	}

	// 浏览量 +1（异步更新，不影响响应）
	go r.incrementViews(article.ID)

	return detail, nil
}

// getArticleTags 获取文章的标签列表
func (r *articlesRepository) getArticleTags(ctx context.Context, articleID int64) []dto.TagItem {
	var tags []dto.TagItem

	// 通过关联表查询标签
	r.db.WithContext(ctx).
		Table("tags").
		Select("tags.id, tags.name, tags.slug").
		Joins("JOIN article_tags ON article_tags.tag_id = tags.id").
		Where("article_tags.article_id = ?", articleID).
		Find(&tags)

	// 确保返回空数组而不是 nil
	if tags == nil {
		tags = []dto.TagItem{}
	}

	return tags
}

// incrementViews 增加浏览量
func (r *articlesRepository) incrementViews(articleID int64) {
	r.db.Model(&entity.Article{}).
		Where("id = ?", articleID).
		UpdateColumn("views", gorm.Expr("views + 1"))
}
