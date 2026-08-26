package repository

import (
	"context"

	"ZZHAN/internal/model/entity"
)

// AuthRepository 认证仓储接口
type AuthRepository interface {
	// FindUserByOpenID 根据 provider 和 openid 查找用户
	FindUserByOpenID(ctx context.Context, provider, openid string) (*entity.User, error)

	// FindUserByID 根据用户ID查找用户
	FindUserByID(ctx context.Context, userID int) (*entity.User, error)

	// CreateUser 创建用户
	CreateUser(ctx context.Context, user *entity.User) error

	// UpdateUser 更新用户信息
	UpdateUser(ctx context.Context, user *entity.User) error

	// UpdateUserLoginTime 更新用户最后登录时间
	UpdateUserLoginTime(ctx context.Context, userID int) error
}
