package service

import (
	"ZZHAN/internal/repository"
	"context"
	"fmt"
	"time"

	"ZZHAN/internal/model/dto"
	"ZZHAN/pkg/config"
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

// accessTTL 获取 access_token 的过期时长
func (s *adminAuthService) accessTTL() time.Duration {
	return config.Get().JWT.AccessExpireHours * time.Hour
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

	// 单设备登录：踢掉旧 token，记录新 token
	if s.redisRepo != nil && s.redisRepo.Available(ctx) {
		adminID := int(admin.ID)
		// 获取旧的活跃 token 并加入黑名单
		if oldToken, err := s.redisRepo.GetActiveToken(ctx, "admin", adminID); err == nil && oldToken != "" {
			_ = s.redisRepo.AddToBlacklist(ctx, oldToken, s.accessTTL())
		}
		// 存储新的活跃 token
		_ = s.redisRepo.SetActiveToken(ctx, "admin", adminID, accessToken, s.accessTTL())
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

// RefreshToken 刷新 access_token
func (s *adminAuthService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	// 验证 refresh_token
	claims, err := jwt.ParseToken(refreshToken, "refresh_token")
	if err != nil {
		return nil, fmt.Errorf("refresh_token 无效：%w", err)
	}

	// 检查是否在黑名单中
	if s.redisRepo != nil && s.redisRepo.Available(ctx) {
		blacklisted, err := s.redisRepo.IsBlacklisted(ctx, refreshToken)
		if err != nil {
			return nil, fmt.Errorf("检查 token 黑名单失败：%w", err)
		}
		if blacklisted {
			return nil, fmt.Errorf("refresh_token 已失效")
		}
	}

	newAccessToken, err := jwt.GenerateToken(claims.UserID, claims.Username, "access_token")
	if err != nil {
		return nil, fmt.Errorf("生成新的 access_token 失败：%w", err)
	}

	// 更新活跃 token（旧的 access_token 已失效，用新的替换）
	if s.redisRepo != nil && s.redisRepo.Available(ctx) {
		_ = s.redisRepo.SetActiveToken(ctx, "admin", claims.UserID, newAccessToken, s.accessTTL())
	}

	return &dto.RefreshResponse{
		AccessToken: newAccessToken,
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
	if err := s.redisRepo.AddToBlacklist(ctx, accessToken, remaining); err != nil {
		return err
	}

	// 清除活跃 token
	if s.redisRepo != nil && s.redisRepo.Available(ctx) {
		_ = s.redisRepo.ClearActiveToken(ctx, "admin", claims.UserID)
	}

	return nil
}

// GetProfile 获取管理员资料
func (s *adminAuthService) GetProfile(ctx context.Context, adminID int) (*dto.AdminProfileResponse, error) {
	admin, err := s.adminRepo.FindAdminByID(ctx, adminID)
	if err != nil {
		return nil, fmt.Errorf("管理员不存在：%w", err)
	}

	return &dto.AdminProfileResponse{
		Username: admin.Username,
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
	if req.Username != "" {
		admin.Username = req.Username
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
		Avatar:   admin.Avatar,
	}, nil
}
