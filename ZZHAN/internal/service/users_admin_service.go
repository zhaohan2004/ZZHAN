package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type usersAdminService struct {
	usersAdminRepo repository.UsersAdminRepository
}

// NewUsersAdminService 创建后台用户业务实例
func NewUsersAdminService(usersAdminRepo repository.UsersAdminRepository) UsersAdminService {
	return &usersAdminService{usersAdminRepo: usersAdminRepo}
}

// AdminList 获取用户列表
func (s *usersAdminService) AdminList(ctx context.Context, req *dto.UsersAdminListRequest) (*dto.UsersAdminListResponse, error) {
	// 分页默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	users, total, err := s.usersAdminRepo.AdminList(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.UsersAdminListResponse{
		List:  users,
		Total: total,
	}, nil
}

// AdminGetByID 获取用户详情
func (s *usersAdminService) AdminGetByID(ctx context.Context, id int64) (*dto.UsersAdminItem, error) {
	return s.usersAdminRepo.AdminGetByID(ctx, id)
}

// AdminUpdateStatus 修改用户状态
func (s *usersAdminService) AdminUpdateStatus(ctx context.Context, id int64, status int8) (*dto.UsersAdminItem, error) {
	if err := s.usersAdminRepo.AdminUpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return s.usersAdminRepo.AdminGetByID(ctx, id)
}

// AdminDelete 删除用户
func (s *usersAdminService) AdminDelete(ctx context.Context, id int64) error {
	return s.usersAdminRepo.AdminDelete(ctx, id)
}
