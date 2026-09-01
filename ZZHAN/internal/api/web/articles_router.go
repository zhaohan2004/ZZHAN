package web

import (
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册文章相关路由
func (c *ArticlesController) RegisterRoutes(r *gin.RouterGroup) {
	// 文章路由组
	articleGroup := r.Group("/articles")
	{
		articleGroup.GET("", c.GetArticleList)                                    // GET /api/v1/articles - 文章列表
		articleGroup.GET("/:slug", middleware.OptionalAuth(), c.GetArticleDetail) // GET /api/v1/articles/:slug - 文章详情（可选认证，用于判断点赞状态）
	}
}
