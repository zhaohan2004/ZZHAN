package repository

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"context"

	"gorm.io/gorm"
)

type commentsRepository struct {
	db *gorm.DB
}

// NewCommentsRepository 创建评论仓储实例
func NewCommentsRepository(db *gorm.DB) CommentsRepository {
	return &commentsRepository{db: db}
}

// GetByArticleSlug 获取文章评论列表（分页）
func (r *commentsRepository) GetByArticleSlug(ctx context.Context, slug string, page, pageSize int) ([]dto.CommentItem, int64, error) {
	// 先通过 slug 查找文章 ID
	var article entity.Article
	if err := r.db.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&article).Error; err != nil {
		return nil, 0, err
	}

	var total int64
	query := r.db.WithContext(ctx).
		Model(&entity.Comment{}).
		Where("article_id = ? AND status = ? AND parent_id IS NULL", article.ID, "normal")

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	var comments []entity.Comment
	offset := (page - 1) * pageSize
	if err := query.
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&comments).Error; err != nil {
		return nil, 0, err
	}

	// 转换为 DTO
	result := make([]dto.CommentItem, 0, len(comments))
	for _, c := range comments {
		item := dto.CommentItem{
			ID:         c.ID,
			ParentID:   c.ParentID,
			UserName:   c.UserName,
			UserAvatar: c.UserAvatar,
			Content:    c.Content,
			CreatedAt:  c.CreatedAt,
			LikeCount:  c.LikeCount,
			Liked:      false,
		}

		// 查询子评论（回复）
		item.Replies = r.getReplies(ctx, c.ID)

		result = append(result, item)
	}

	return result, total, nil
}

// getReplies 获取评论的回复列表
func (r *commentsRepository) getReplies(ctx context.Context, parentID int64) []dto.CommentItem {
	var replies []entity.Comment
	r.db.WithContext(ctx).
		Where("parent_id = ? AND status = ?", parentID, "normal").
		Order("created_at ASC").
		Find(&replies)

	result := make([]dto.CommentItem, 0, len(replies))
	for _, reply := range replies {
		result = append(result, dto.CommentItem{
			ID:         reply.ID,
			ParentID:   reply.ParentID,
			UserName:   reply.UserName,
			UserAvatar: reply.UserAvatar,
			Content:    reply.Content,
			CreatedAt:  reply.CreatedAt,
			LikeCount:  reply.LikeCount,
			Liked:      false,
		})
	}

	if result == nil {
		result = []dto.CommentItem{}
	}

	return result
}

// Create 创建评论
func (r *commentsRepository) Create(ctx context.Context, articleID, userID int64, userName, userAvatar, content, ip string, parentID *int64) (int64, error) {
	comment := entity.Comment{
		ArticleID:  articleID,
		ParentID:   parentID,
		UserID:     userID,
		UserName:   userName,
		UserAvatar: userAvatar,
		Content:    content,
		IP:         ip,
		Status:     "normal",
		LikeCount:  0,
	}

	if err := r.db.WithContext(ctx).Create(&comment).Error; err != nil {
		return 0, err
	}

	// 更新文章评论数
	r.db.WithContext(ctx).
		Model(&entity.Article{}).
		Where("id = ?", articleID).
		UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	return comment.ID, nil
}

// GetArticleIDBySlug 通过 slug 获取文章 ID
func (r *commentsRepository) GetArticleIDBySlug(ctx context.Context, slug string) (int64, error) {
	var article entity.Article
	if err := r.db.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&article).Error; err != nil {
		return 0, err
	}
	return article.ID, nil
}

// GetUserAvatarByID 通过用户ID获取头像
func (r *commentsRepository) GetUserAvatarByID(ctx context.Context, userID int64) (string, error) {
	var user entity.User
	if err := r.db.WithContext(ctx).Select("avatar").Where("id = ?", userID).First(&user).Error; err != nil {
		return "", err
	}
	return user.Avatar, nil
}
