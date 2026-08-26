package api

import (
	"ZZHAN/internal/api/admin"
	"ZZHAN/internal/api/web"
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 路由
type Router struct {
	siteController      *web.SiteController
	siteAdminController *admin.SiteAdminController
	adminAuthController *admin.AdminAuthController
	authController      *web.AuthController
}

// NewRouter 创建路由
func NewRouter(
	siteController *web.SiteController,
	siteAdminController *admin.SiteAdminController,
	adminAuthController *admin.AdminAuthController,
	authController *web.AuthController,
) *Router {
	return &Router{
		siteController:      siteController,
		siteAdminController: siteAdminController,
		adminAuthController: adminAuthController,
		authController:      authController,
	}
}

// Setup 设置路由
func (r *Router) Setup(engine *gin.Engine) {
	// 全局中间件
	engine.Use(middleware.Recovery())
	engine.Use(middleware.RequestLogger())
	engine.Use(middleware.CORS())

	// 健康检查
	engine.GET("/api/v1/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "ZZHAN API is running",
		})
	})

	// API 路由组
	apiGroup := engine.Group("/api/v1")
	{
		_ = apiGroup
		// ========== 在这里注册你的路由 ==========
		//前台获取站点信息路由
		r.siteController.RegisterRoutes(apiGroup)

		//后台认证路由（验证码 + 登录）
		r.adminAuthController.RegisterRoutes(apiGroup)

		//后台系统设置路由
		r.siteAdminController.RegisterRoutes(apiGroup)

		//前台认证路由
		r.authController.RegisterRoutes(apiGroup)
	}
}

// Close 关闭所有路由连接
func (r *Router) Close() error {
	return nil
}
