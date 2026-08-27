package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CommentsService 评论业务接口
type CommentsService interface {
	// GetByArticleSlug 获取文章评论列表
	GetByArticleSlug(ctx context.Context, slug string, page, pageSize int) (*dto.CommentListResponse, error)

	// Create 创建评论
	Create(ctx context.Context, slug string, userID int64, userName, userAvatar, ip string, req *dto.CommentCreateRequest) (*dto.CommentCreateResult, error)
}
