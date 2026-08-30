package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// UsersAdminService 后台用户业务接口
type UsersAdminService interface {
	// AdminList 获取用户列表
	AdminList(ctx context.Context, req *dto.UsersAdminListRequest) (*dto.UsersAdminListResponse, error)

	// AdminGetByID 获取用户详情
	AdminGetByID(ctx context.Context, id int64) (*dto.UsersAdminItem, error)

	// AdminUpdateStatus 修改用户状态
	AdminUpdateStatus(ctx context.Context, id int64, status int8) (*dto.UsersAdminItem, error)

	// AdminDelete 删除用户
	AdminDelete(ctx context.Context, id int64) error
}
