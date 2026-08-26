package repository

import (
	"context"

	"ZZHAN/internal/model/entity"
)

// AdminAuthRepository 后台认证仓储接口
type AdminAuthRepository interface {
	// FindAdminByUsername 根据用户名查找管理员
	FindAdminByUsername(ctx context.Context, username string) (*entity.Admin, error)

	// FindAdminByID 根据 ID 查找管理员
	FindAdminByID(ctx context.Context, id int) (*entity.Admin, error)

	// UpdateAdmin 更新管理员信息
	UpdateAdmin(ctx context.Context, admin *entity.Admin) error
}
