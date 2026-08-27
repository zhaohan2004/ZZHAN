package repository

import "context"

// LikeRepository 点赞仓储接口
type LikeRepository interface {
	// ToggleArticleLike 文章点赞，返回是否点赞和最新点赞数
	ToggleArticleLike(ctx context.Context, articleID, userID int64) (liked bool, likes int32, err error)

	// ToggleCommentLike 评论点赞，返回是否点赞和最新点赞数
	ToggleCommentLike(ctx context.Context, commentID, userID int64) (liked bool, likeCount int32, err error)

	// GetArticleIDBySlug 通过 slug 获取文章 ID
	GetArticleIDBySlug(ctx context.Context, slug string) (int64, error)
}
