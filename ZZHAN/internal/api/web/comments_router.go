package web

import (
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册评论相关路由
func (c *CommentsController) RegisterRoutes(r *gin.RouterGroup) {
	comments := r.Group("/comments")
	{
		comments.GET("/:slug", c.GetComments) // GET /api/v1/comments/:slug - 评论列表（公开）

		// 发表评论需要登录
		auth := comments.Group("")
		auth.Use(middleware.Auth(c.redisRepo))
		{
			auth.POST("/:slug", c.PostComment) // POST /api/v1/comments/:slug - 发表评论（需登录）
		}
	}
}
