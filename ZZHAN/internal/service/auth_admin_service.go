package service

import (
	"ZZHAN/internal/repository"
	"context"
	"fmt"
	"time"

	"ZZHAN/internal/model/dto"
	appErrors "ZZHAN/pkg/errors"
	"ZZHAN/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type adminAuthService struct {
	adminRepo repository.AdminAuthRepository
	redisRepo repository.RedisRepository
}

// NewAdminAuthService 创建后台认证业务
func NewAdminAuthService(adminRepo repository.AdminAuthRepository, redisRepo repository.RedisRepository) AdminAuthService {
	return &adminAuthService{adminRepo: adminRepo, redisRepo: redisRepo}
}

// Login 后台管理员登录
func (s *adminAuthService) Login(ctx context.Context, req *dto.LoginAdminRequest) (*dto.LoginResponse, error) {
	// 查找管理员
	admin, err := s.adminRepo.FindAdminByUsername(ctx, req.UserName)
	if err != nil {
		return nil, appErrors.ErrInvalidCredentials
	}

	// 校验密码
	if err := bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(req.Password)); err != nil {
		return nil, appErrors.ErrInvalidCredentials
	}

	// 生成 JWT Token 对
	accessToken, refreshToken, err := jwt.GenerateTokenPair(int(admin.ID), admin.Username)
	if err != nil {
		return nil, appErrors.NewWithErr(appErrors.CodeInternalError, "生成令牌失败", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.AuthUser{
			ID:       int(admin.ID),
			Provider: "admin",
			Nickname: admin.Nickname,
			Avatar:   admin.Avatar,
		},
	}, nil
}

// Logout 退出登录
func (s *adminAuthService) Logout(ctx context.Context, accessToken string) error {
	claims, err := jwt.ParseToken(accessToken, "access_token")
	if err != nil {
		// token 已失效，直接返回成功
		return nil
	}

	// 计算剩余过期时间
	expireAt := claims.ExpiresAt.Time
	remaining := time.Until(expireAt)
	if remaining <= 0 {
		return nil
	}

	// 加入黑名单
	return s.redisRepo.AddToBlacklist(ctx, accessToken, remaining)
}

// GetProfile 获取管理员资料
func (s *adminAuthService) GetProfile(ctx context.Context, adminID int) (*dto.AdminProfileResponse, error) {
	admin, err := s.adminRepo.FindAdminByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("管理员不存在：%w", err)
	}

	return &dto.AdminProfileResponse{
		Username: admin.Username,
		Nickname: admin.Nickname,
		Avatar:   admin.Avatar,
	}, nil
}

// UpdateProfile 更新管理员资料
func (s *adminAuthService) UpdateProfile(ctx context.Context, adminID int, req *dto.UpdateAdminProfileRequest) (*dto.AdminProfileResponse, error) {
	admin, err := s.adminRepo.FindAdminByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("管理员不存在：%w", err)
	}

	// 更新资料字段
	if req.Nickname != "" {
		admin.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		admin.Avatar = req.Avatar
	}

	// 如果传了新密码，重新哈希
	if req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, appErrors.NewWithErr(appErrors.CodeInternalError, "密码加密失败", err)
		}
		admin.PasswordHash = string(hash)
	}

	if err := s.adminRepo.UpdateAdmin(ctx, admin); err != nil {
		return nil, fmt.Errorf("更新管理员信息失败：%w", err)
	}

	return &dto.AdminProfileResponse{
		Username: admin.Username,
		Nickname: admin.Nickname,
		Avatar:   admin.Avatar,
	}, nil
}
