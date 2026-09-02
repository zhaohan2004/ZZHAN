package middleware

import (
	"ZZHAN/internal/repository"
	"strings"

	"ZZHAN/pkg/jwt"
	"ZZHAN/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	// ContextUserID 用户 ID 上下文键
	ContextUserID = "user_id"
	// ContextUsername 用户名 上下文键
	ContextUsername = "username"
)

// Auth JWT 认证中间件
// userType: "admin" 或 "user"，用于单设备登录校验；为空则跳过活跃 token 检查
func Auth(redisRepo repository.RedisRepository, userType ...string) gin.HandlerFunc {
	ut := ""
	if len(userType) > 0 {
		ut = userType[0]
	}

	return func(c *gin.Context) {
		// 从 Header 获取 Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "请提供认证令牌")
			c.Abort()
			return
		}

		// 解析 Bearer Token
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "令牌格式错误")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 解析 Token
		claims, err := jwt.ParseToken(tokenString, "access_token")
		if err != nil {
			response.Unauthorized(c, err.Error())
			c.Abort()
			return
		}

		// 检查 token 是否存在于黑名单中
		if redisRepo != nil && redisRepo.Available(c) {
			blacklisted, err := redisRepo.IsBlacklisted(c, tokenString)
			if err != nil {
				response.Unauthorized(c, "搜索黑名单失败")
				c.Abort()
				return
			}
			if blacklisted {
				response.Unauthorized(c, "token 已失效")
				c.Abort()
				return
			}

			// 单设备登录校验：检查 token 是否为当前用户的活跃 token
			if ut != "" {
				activeToken, err := redisRepo.GetActiveToken(c, ut, claims.GetUserID())
				if err == nil && activeToken != "" && activeToken != tokenString {
					response.Unauthorized(c, "账号已在其他设备登录")
					c.Abort()
					return
				}
			}
		}

		// 将用户信息存入上下文
		c.Set(ContextUserID, claims.GetUserID())
		c.Set(ContextUsername, claims.GetUsername())

		c.Next()
	}
}

// GetUserID 从上下文获取用户 ID
func GetUserID(c *gin.Context) int {
	if userID, exists := c.Get(ContextUserID); exists {
		return userID.(int)
	}
	return 0
}

// GetUsername 从上下文获取用户名
func GetUsername(c *gin.Context) string {
	if username, exists := c.Get(ContextUsername); exists {
		return username.(string)
	}
	return ""
}

// OptionalAuth 可选的 JWT 认证中间件
func OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.Next()
			return
		}

		claims, err := jwt.ParseToken(parts[1], "access_token")
		if err == nil {
			c.Set(ContextUserID, claims.GetUserID())
			c.Set(ContextUsername, claims.GetUsername())
		}

		c.Next()
	}
}
