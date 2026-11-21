package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	appErrors "ZZHAN/pkg/errors"
	"ZZHAN/pkg/jwt"
	"ZZHAN/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	// ContextUserID 用户 ID 上下文键
	ContextUserID = "user_id"
	// ContextUsername 用户名 上下文键
	ContextUsername = "username"
)

// Auth JWT 认证中间件
// userType 用于单设备登录校验，为空时跳过活跃 token 检查。
func Auth(redisRepo repository.RedisRepository, userType ...string) gin.HandlerFunc {
	scope := ""
	if len(userType) > 0 {
		scope = userType[0]
	}
	return authHandler(redisRepo, func(context.Context, *jwt.CustomClaims) (string, error) {
		return scope, nil
	}, false)
}

// AuthUser requires an enabled ordinary user with a current session.
func AuthUser(redisRepo repository.RedisRepository, authRepo repository.AuthRepository) gin.HandlerFunc {
	return userAuth(redisRepo, authRepo, false)
}

// OptionalUserAuth treats requests without a current user session as anonymous.
func OptionalUserAuth(redisRepo repository.RedisRepository, authRepo repository.AuthRepository) gin.HandlerFunc {
	return userAuth(redisRepo, authRepo, true)
}

func userAuth(redisRepo repository.RedisRepository, authRepo repository.AuthRepository, optional bool) gin.HandlerFunc {
	return authHandler(redisRepo, func(ctx context.Context, claims *jwt.CustomClaims) (string, error) {
		if err := service.ValidateUserSession(ctx, authRepo, claims.UserID, claims.SessionVersion); err != nil {
			return "", err
		}
		return service.UserSessionScope(claims.SessionVersion), nil
	}, optional)
}

func authHandler(redisRepo repository.RedisRepository, sessionScope func(context.Context, *jwt.CustomClaims) (string, error), optional bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		claims, err := authenticate(c.Request.Context(), c.GetHeader("Authorization"), redisRepo, sessionScope)
		if err != nil {
			if optional {
				c.Next()
				return
			}
			switch {
			case errors.Is(err, appErrors.ErrUserDisabled):
				response.Error(c, http.StatusForbidden, appErrors.CodeUserDisabled, appErrors.ErrUserDisabled.Message)
			case errors.Is(err, appErrors.ErrInternalError):
				_ = c.Error(err)
				response.Error(c, http.StatusServiceUnavailable, appErrors.CodeInternalError, "暂时无法验证登录状态")
			case errors.Is(err, appErrors.ErrInvalidToken):
				response.Unauthorized(c, "登录已失效，请重新登录")
			default:
				response.Unauthorized(c, err.Error())
			}
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.GetUserID())
		c.Set(ContextUsername, claims.GetUsername())
		c.Next()
	}
}

func authenticate(ctx context.Context, header string, redisRepo repository.RedisRepository, sessionScope func(context.Context, *jwt.CustomClaims) (string, error)) (*jwt.CustomClaims, error) {
	if header == "" {
		return nil, errors.New("请提供认证令牌")
	}
	scheme, token, found := strings.Cut(header, " ")
	if !found || scheme != "Bearer" {
		return nil, errors.New("令牌格式错误")
	}
	claims, err := jwt.ParseToken(token, "access_token")
	if err != nil {
		return nil, err
	}
	scope, err := sessionScope(ctx, claims)
	if err != nil {
		return nil, err
	}
	if err := checkActiveToken(ctx, redisRepo, scope, claims.UserID, token); err != nil {
		return nil, err
	}
	return claims, nil
}

func checkActiveToken(ctx context.Context, redisRepo repository.RedisRepository, scope string, userID int, token string) error {
	if redisRepo == nil || !redisRepo.Available(ctx) {
		return nil
	}
	blacklisted, err := redisRepo.IsBlacklisted(ctx, token)
	if err != nil {
		return errors.New("搜索黑名单失败")
	}
	if blacklisted {
		return errors.New("token 已失效")
	}
	if scope == "" {
		return nil
	}
	activeToken, err := redisRepo.GetActiveToken(ctx, scope, userID)
	if errors.Is(err, redis.Nil) {
		return errors.New("登录已失效，请重新登录")
	}
	if err == nil && activeToken != token {
		return errors.New("账号已在其他设备登录")
	}
	return nil
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
