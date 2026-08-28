package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type commentsService struct {
	commentsRepo repository.CommentsRepository
}

// NewCommentsService 创建评论业务实例
func NewCommentsService(commentsRepo repository.CommentsRepository) CommentsService {
	return &commentsService{commentsRepo: commentsRepo}
}

// GetByArticleSlug 获取文章评论列表
func (s *commentsService) GetByArticleSlug(ctx context.Context, slug string, page, pageSize int) (*dto.CommentListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := s.commentsRepo.GetByArticleSlug(ctx, slug, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.CommentListResponse{
		List:  list,
		Total: total,
	}, nil
}

// GetReplies 获取评论的回复列表
func (s *commentsService) GetReplies(ctx context.Context, commentID int64, page, pageSize int) (*dto.CommentListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}

	list, total, err := s.commentsRepo.GetReplies(ctx, commentID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &dto.CommentListResponse{
		List:  list,
		Total: total,
	}, nil
}

// Create 创建评论
func (s *commentsService) Create(ctx context.Context, slug string, userID int64, userName, userAvatar, ip string, req *dto.CommentCreateRequest) (*dto.CommentCreateResult, error) {
	// 需要先通过 slug 获取文章 ID（复用 repository 的能力）
	// 这里直接传 slug 给 repository，让它内部查文章
	// 但 repository 的 Create 接口需要 articleID，所以需要调整
	// 先用一个辅助方法获取 articleID
	articleID, err := s.commentsRepo.GetArticleIDBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	id, err := s.commentsRepo.Create(ctx, articleID, userID, userName, userAvatar, req.Content, ip, req.ParentID)
	if err != nil {
		return nil, err
	}

	return &dto.CommentCreateResult{
		ID:      id,
		Status:  "normal",
		Message: "评论已提交",
	}, nil
}
