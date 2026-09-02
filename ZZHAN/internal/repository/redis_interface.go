package repository

import (
	"context"
	"time"
)

// RedisRepository Redis 仓储接口
type RedisRepository interface {
	// Set 设置缓存
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// Get 获取缓存
	Get(ctx context.Context, key string, dest interface{}) error
	// Del 删除缓存
	Del(ctx context.Context, key string) error
	// Exists 检查是否存在
	Exists(ctx context.Context, key string) (bool, error)
	// Available 检查redis是否可用
	Available(ctx context.Context) bool
	// ========== Token 黑名单方法 ==========

	// AddToBlacklist 将 token 加入黑名单
	AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error

	// IsBlacklisted 检查 token 是否在黑名单中
	IsBlacklisted(ctx context.Context, token string) (bool, error)

	// ========== 活跃 Token 管理（单设备登录） ==========

	// SetActiveToken 设置用户当前活跃的 access_token（登录/刷新时调用）
	SetActiveToken(ctx context.Context, userType string, userID int, token string, expiration time.Duration) error

	// GetActiveToken 获取用户当前活跃的 access_token
	GetActiveToken(ctx context.Context, userType string, userID int) (string, error)

	// SetActiveRefreshToken 设置用户当前活跃的 refresh_token
	SetActiveRefreshToken(ctx context.Context, userType string, userID int, token string, expiration time.Duration) error

	// GetActiveRefreshToken 获取用户当前活跃的 refresh_token
	GetActiveRefreshToken(ctx context.Context, userType string, userID int) (string, error)

	// ClearActiveToken 清除用户的活跃 token（退出时调用，同时清除 access 和 refresh）
	ClearActiveToken(ctx context.Context, userType string, userID int) error

	// ========== 浏览量去重方法 ==========

	// CheckViewAccess 检查是否已访问过文章（用于浏览量去重）
	// 返回 true 表示已访问过，false 表示首次访问
	CheckViewAccess(ctx context.Context, articleID int64, clientIP string) (bool, error)

	// SetViewAccess 设置文章访问标记
	SetViewAccess(ctx context.Context, articleID int64, clientIP string, expiration time.Duration) error
}
