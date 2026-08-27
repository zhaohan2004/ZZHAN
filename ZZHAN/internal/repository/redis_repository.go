package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisRepository struct {
	redis *redis.Client
}

// NewRedisRepository 创建 Redis 仓库
func NewRedisRepository(redisClient *redis.Client) RedisRepository {
	return &redisRepository{
		redis: redisClient,
	}
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if r.redis == nil {
		return nil
	}
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redis.Set(ctx, key, data, expiration).Err()
}

func (r *redisRepository) Get(ctx context.Context, key string, dest interface{}) error {
	if r.redis == nil {
		return redis.Nil
	}
	data, err := r.redis.Get(ctx, key).Bytes()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	if r.redis == nil {
		return nil
	}
	return r.redis.Del(ctx, key).Err()
}

func (r *redisRepository) Exists(ctx context.Context, key string) (bool, error) {
	if r.redis == nil {
		return false, nil
	}
	n, err := r.redis.Exists(ctx, key).Result()
	return n > 0, err
}

// Available 判断 Redis 是否可用
func (r *redisRepository) Available(ctx context.Context) bool {
	if r.redis == nil {
		return false
	}
	return r.redis.Ping(ctx).Err() == nil
}

// ========== Token 黑名单方法 ==========

// blacklistKeyPrefix 黑名单 key 前缀
const blacklistKeyPrefix = "token:blacklist:"

// AddToBlacklist 将 token 加入黑名单
func (r *redisRepository) AddToBlacklist(ctx context.Context, token string, expiration time.Duration) error {
	if r.redis == nil {
		return nil
	}
	key := blacklistKeyPrefix + token
	return r.redis.Set(ctx, key, "1", expiration).Err()
}

// IsBlacklisted 检查 token 是否在黑名单中
func (r *redisRepository) IsBlacklisted(ctx context.Context, token string) (bool, error) {
	if r.redis == nil {
		return false, nil
	}
	key := blacklistKeyPrefix + token
	n, err := r.redis.Exists(ctx, key).Result()
	return n > 0, err
}

// ========== 浏览量去重方法 ==========

// viewAccessKeyPrefix 浏览量去重 key 前缀
const viewAccessKeyPrefix = "article:views:"

// CheckViewAccess 检查是否已访问过文章
func (r *redisRepository) CheckViewAccess(ctx context.Context, articleID int64, clientIP string) (bool, error) {
	if r.redis == nil {
		return false, nil
	}
	key := fmt.Sprintf("%s%d:%s", viewAccessKeyPrefix, articleID, clientIP)
	n, err := r.redis.Exists(ctx, key).Result()
	return n > 0, err
}

// SetViewAccess 设置文章访问标记
func (r *redisRepository) SetViewAccess(ctx context.Context, articleID int64, clientIP string, expiration time.Duration) error {
	if r.redis == nil {
		return nil
	}
	key := fmt.Sprintf("%s%d:%s", viewAccessKeyPrefix, articleID, clientIP)
	return r.redis.Set(ctx, key, "1", expiration).Err()
}
