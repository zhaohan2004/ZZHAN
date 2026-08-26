package web

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册认证相关路由
func (c *AuthController) RegisterRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	auth.Use(middleware.Auth(c.redisRepo))
	{
		auth.GET("/me", c.GetCurrentUser)
		auth.POST("/refresh", c.RefreshToken)
		auth.POST("/logout", c.Logout)
	}

	r.POST("/auth/github", c.GitHubLogin)
}
