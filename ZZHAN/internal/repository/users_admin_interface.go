package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// UsersAdminRepository 后台用户仓储接口
type UsersAdminRepository interface {
	// AdminList 获取用户列表（支持筛选）
	AdminList(ctx context.Context, req *dto.UsersAdminListRequest) ([]dto.UsersAdminItem, int64, error)

	// AdminGetByID 通过 ID 获取用户详情
	AdminGetByID(ctx context.Context, id int64) (*dto.UsersAdminItem, error)

	// AdminUpdateStatus 修改用户状态
	AdminUpdateStatus(ctx context.Context, id int64, status int8) error

	// AdminDelete 删除用户
	AdminDelete(ctx context.Context, id int64) error
}
