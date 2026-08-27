package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type likeService struct {
	likeRepo repository.LikeRepository
}

// NewLikeService 创建点赞业务实例
func NewLikeService(likeRepo repository.LikeRepository) LikeService {
	return &likeService{likeRepo: likeRepo}
}

// ArticleLike 文章点赞
func (s *likeService) ArticleLike(ctx context.Context, slug string, userID int64) (*dto.ArticleLikeResult, error) {
	articleID, err := s.likeRepo.GetArticleIDBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}

	liked, likes, err := s.likeRepo.ToggleArticleLike(ctx, articleID, userID)
	if err != nil {
		return nil, err
	}

	return &dto.ArticleLikeResult{
		Liked: liked,
		Likes: likes,
	}, nil
}

// CommentLike 评论点赞
func (s *likeService) CommentLike(ctx context.Context, commentID, userID int64) (*dto.CommentLikeResult, error) {
	liked, likeCount, err := s.likeRepo.ToggleCommentLike(ctx, commentID, userID)
	if err != nil {
		return nil, err
	}

	return &dto.CommentLikeResult{
		Liked:     liked,
		LikeCount: likeCount,
	}, nil
}
