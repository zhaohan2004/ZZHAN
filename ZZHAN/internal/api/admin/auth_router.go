package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *AdminAuthController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	{
		// 公开接口：验证码 + 登录 + 刷新 token
		admin.GET("/auth/captcha", c.GetCaptcha)
		admin.POST("/auth/login", c.Login)
		admin.POST("/auth/refresh", c.RefreshToken)

		// 需要登录态的接口（经过 Auth 中间件校验 token + 黑名单）
		auth := admin.Group("")
		auth.Use(middleware.Auth(c.redisRepo))
		{
			auth.POST("/auth/logout", c.Logout)
			auth.GET("/profile", c.GetProfile)
			auth.PUT("/profile", c.UpdateProfile)
		}
	}
}
