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
}
