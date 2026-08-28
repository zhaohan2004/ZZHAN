package api

import (
	"ZZHAN/internal/api/admin"
	"ZZHAN/internal/api/web"
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

// Router 路由
type Router struct {
	siteController            *web.SiteController
	siteAdminController       *admin.SiteAdminController
	adminAuthController       *admin.AdminAuthController
	authController            *web.AuthController
	articlesController        *web.ArticlesController
	adminArticlesController   *admin.AdminArticlesController
	categoriesController      *web.CategoriesController
	adminCategoriesController *admin.AdminCategoriesController
	tagsController            *web.TagsController
	adminTagsController       *admin.AdminTagsController
	archivesController        *web.ArchivesController
	statsController           *web.StatsController
	aboutController           *web.AboutController
	commentsController        *web.CommentsController
	likeController            *web.LikeController
	uploadController          *admin.UploadController
}

// NewRouter 创建路由
func NewRouter(
	siteController *web.SiteController,
	siteAdminController *admin.SiteAdminController,
	adminAuthController *admin.AdminAuthController,
	authController *web.AuthController,
	articlesController *web.ArticlesController,
	adminArticlesController *admin.AdminArticlesController,
	categoriesController *web.CategoriesController,
	adminCategoriesController *admin.AdminCategoriesController,
	tagsController *web.TagsController,
	adminTagsController *admin.AdminTagsController,
	archivesController *web.ArchivesController,
	statsController *web.StatsController,
	aboutController *web.AboutController,
	commentsController *web.CommentsController,
	likeController *web.LikeController,
	uploadController *admin.UploadController,
) *Router {
	return &Router{
		siteController:            siteController,
		siteAdminController:       siteAdminController,
		adminAuthController:       adminAuthController,
		authController:            authController,
		articlesController:        articlesController,
		adminArticlesController:   adminArticlesController,
		categoriesController:      categoriesController,
		adminCategoriesController: adminCategoriesController,
		tagsController:            tagsController,
		adminTagsController:       adminTagsController,
		archivesController:        archivesController,
		statsController:           statsController,
		aboutController:           aboutController,
		commentsController:        commentsController,
		likeController:            likeController,
		uploadController:          uploadController,
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

		//前台文章路由
		r.articlesController.RegisterRoutes(apiGroup)

		//后台文章路由
		r.adminArticlesController.RegisterRoutes(apiGroup)

		//前台分类路由
		r.categoriesController.RegisterRoutes(apiGroup)

		//后台分类路由
		r.adminCategoriesController.RegisterRoutes(apiGroup)

		//前台标签路由
		r.tagsController.RegisterRoutes(apiGroup)

		//后台标签路由
		r.adminTagsController.RegisterRoutes(apiGroup)

		//前台归档路由
		r.archivesController.RegisterRoutes(apiGroup)

		//前台统计路由
		r.statsController.RegisterRoutes(apiGroup)

		//前台关于我路由
		r.aboutController.RegisterRoutes(apiGroup)

		//前台评论路由
		r.commentsController.RegisterRoutes(apiGroup)

		//前台点赞路由
		r.likeController.RegisterRoutes(apiGroup)

		//后台上传路由
		r.uploadController.RegisterRoutes(apiGroup)
	}
}

// Close 关闭所有路由连接
func (r *Router) Close() error {
	return nil
}
