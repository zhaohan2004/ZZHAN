package service

import (
	"context"

	"ZZHAN/internal/model/dto"
)

// AuthService 认证业务接口
type AuthService interface {
	// GitHubLogin GitHub OAuth 登录
	GitHubLogin(ctx context.Context, code, redirectURI string) (*dto.LoginResponse, error)

	// RefreshToken 刷新 access_token
	RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error)

	// Logout 退出登录（将 token 加入黑名单）
	Logout(ctx context.Context, accessToken string) error

	// GetCurrentUser 获取当前登录用户信息
	GetCurrentUser(ctx context.Context, userID int) (*dto.AuthUser, error)
}
