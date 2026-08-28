package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CommentsRepository 评论仓储接口
type CommentsRepository interface {
	// GetByArticleSlug 获取文章评论列表（分页）
	GetByArticleSlug(ctx context.Context, slug string, page, pageSize int) ([]dto.CommentItem, int64, error)

	// GetReplies 获取评论的回复列表（分页）
	GetReplies(ctx context.Context, commentID int64, page, pageSize int) ([]dto.CommentItem, int64, error)

	// Create 创建评论
	Create(ctx context.Context, articleID, userID int64, userName, userAvatar, content, ip string, parentID *int64) (int64, error)

	// GetArticleIDBySlug 通过 slug 获取文章 ID
	GetArticleIDBySlug(ctx context.Context, slug string) (int64, error)

	// GetUserAvatarByID 通过用户ID获取头像
	GetUserAvatarByID(ctx context.Context, userID int64) (string, error)
}
