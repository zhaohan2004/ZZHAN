package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// LikeService 点赞业务接口
type LikeService interface {
	// ArticleLike 文章点赞（切换状态）
	ArticleLike(ctx context.Context, slug string, userID int64) (*dto.ArticleLikeResult, error)

	// CommentLike 评论点赞（切换状态）
	CommentLike(ctx context.Context, commentID, userID int64) (*dto.CommentLikeResult, error)
}
