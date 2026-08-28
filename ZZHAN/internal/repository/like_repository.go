package repository

import (
	"ZZHAN/internal/model/entity"
	"context"
	"errors"

	"gorm.io/gorm"
)

type likeRepository struct {
	db *gorm.DB
}

// NewLikeRepository 创建点赞仓储实例
func NewLikeRepository(db *gorm.DB) LikeRepository {
	return &likeRepository{db: db}
}

// ToggleArticleLike 文章点赞
func (r *likeRepository) ToggleArticleLike(ctx context.Context, articleID, userID int64) (bool, int32, error) {
	var like entity.Like
	err := r.db.WithContext(ctx).
		Where("article_id = ? AND user_id = ?", articleID, userID).
		First(&like).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// 未点赞 → 插入
		like = entity.Like{ArticleID: articleID, UserID: userID}
		if err := r.db.WithContext(ctx).Create(&like).Error; err != nil {
			return false, 0, err
		}
		// 计数+1
		r.db.WithContext(ctx).
			Model(&entity.Article{}).
			Where("id = ?", articleID).
			UpdateColumn("likes", gorm.Expr("likes + 1"))
		// 获取最新计数
		var count int32
		r.db.WithContext(ctx).Model(&entity.Article{}).Where("id = ?", articleID).Pluck("likes", &count)
		return true, count, nil
	}

	if err != nil {
		return false, 0, err
	}

	// 已点赞 → 删除
	r.db.WithContext(ctx).Delete(&like)
	// 计数-1
	r.db.WithContext(ctx).
		Model(&entity.Article{}).
		Where("id = ? AND likes > 0", articleID).
		UpdateColumn("likes", gorm.Expr("likes - 1"))
	// 获取最新计数
	var count int32
	r.db.WithContext(ctx).Model(&entity.Article{}).Where("id = ?", articleID).Pluck("likes", &count)
	return false, count, nil
}

// ToggleCommentLike 评论点赞
func (r *likeRepository) ToggleCommentLike(ctx context.Context, commentID, userID int64) (bool, int32, error) {
	var like entity.CommentLike
	err := r.db.WithContext(ctx).
		Where("comment_id = ? AND user_id = ?", commentID, userID).
		First(&like).Error

	if err == gorm.ErrRecordNotFound {
		// 未点赞 → 插入
		like = entity.CommentLike{CommentID: commentID, UserID: userID}
		if err := r.db.WithContext(ctx).Create(&like).Error; err != nil {
			return false, 0, err
		}
		// 计数+1
		r.db.WithContext(ctx).
			Model(&entity.Comment{}).
			Where("id = ?", commentID).
			UpdateColumn("like_count", gorm.Expr("like_count + 1"))
		// 获取最新计数
		var count int32
		r.db.WithContext(ctx).Model(&entity.Comment{}).Where("id = ?", commentID).Pluck("like_count", &count)
		return true, count, nil
	}

	if err != nil {
		return false, 0, err
	}

	// 已点赞 → 删除
	r.db.WithContext(ctx).Delete(&like)
	// 计数-1
	r.db.WithContext(ctx).
		Model(&entity.Comment{}).
		Where("id = ? AND like_count > 0", commentID).
		UpdateColumn("like_count", gorm.Expr("like_count - 1"))
	// 获取最新计数
	var count int32
	r.db.WithContext(ctx).Model(&entity.Comment{}).Where("id = ?", commentID).Pluck("like_count", &count)
	return false, count, nil
}

// GetArticleIDBySlug 通过 slug 获取文章 ID
func (r *likeRepository) GetArticleIDBySlug(ctx context.Context, slug string) (int64, error) {
	var article entity.Article
	if err := r.db.WithContext(ctx).Select("id").Where("slug = ?", slug).First(&article).Error; err != nil {
		return 0, err
	}
	return article.ID, nil
}
