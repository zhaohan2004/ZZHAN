package service

import (
	"context"

	"ZZHAN/internal/model/dto"
)

// AdminAuthService 后台认证业务接口
type AdminAuthService interface {
	// Login 后台管理员登录
	Login(ctx context.Context, req *dto.LoginAdminRequest) (*dto.LoginResponse, error)

	// Logout 退出登录（将 token 加入黑名单）
	Logout(ctx context.Context, accessToken string) error

	// GetProfile 获取管理员资料
	GetProfile(ctx context.Context, adminID int) (*dto.AdminProfileResponse, error)

	// UpdateProfile 更新管理员资料
	UpdateProfile(ctx context.Context, adminID int, req *dto.UpdateAdminProfileRequest) (*dto.AdminProfileResponse, error)
}
