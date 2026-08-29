package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type commentsAdminService struct {
	commentsAdminRepo repository.CommentsAdminRepository
}

// NewCommentsAdminService 创建后台评论业务实例
func NewCommentsAdminService(commentsAdminRepo repository.CommentsAdminRepository) CommentsAdminService {
	return &commentsAdminService{commentsAdminRepo: commentsAdminRepo}
}

// AdminList 获取评论列表
func (s *commentsAdminService) AdminList(ctx context.Context, req *dto.CommentsAdminListRequest) (*dto.CommentsAdminListResponse, error) {
	// 分页默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	comments, total, err := s.commentsAdminRepo.AdminList(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.CommentsAdminListResponse{
		List:  comments,
		Total: total,
	}, nil
}

// AdminGetByID 获取评论详情
func (s *commentsAdminService) AdminGetByID(ctx context.Context, id int64) (*dto.CommentsAdminItem, error) {
	return s.commentsAdminRepo.AdminGetByID(ctx, id)
}

// AdminUpdateStatus 修改评论状态
func (s *commentsAdminService) AdminUpdateStatus(ctx context.Context, id int64, status string) (*dto.CommentsAdminItem, error) {
	if err := s.commentsAdminRepo.AdminUpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return s.commentsAdminRepo.AdminGetByID(ctx, id)
}

// AdminDelete 删除评论
func (s *commentsAdminService) AdminDelete(ctx context.Context, id int64) error {
	return s.commentsAdminRepo.AdminDelete(ctx, id)
}
