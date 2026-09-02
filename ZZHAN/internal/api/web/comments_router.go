package web

import (
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册评论相关路由
func (c *CommentsController) RegisterRoutes(r *gin.RouterGroup) {
	comments := r.Group("/comments")
	{
		comments.GET("/:slug", c.GetComments)      // 评论列表（公开）
		comments.GET("/replies/:id", c.GetReplies) // 回复列表（公开）

		// 发表评论需要登录
		auth := comments.Group("")
		auth.Use(middleware.Auth(c.redisRepo, "user"))
		{
			auth.POST("/:slug", c.PostComment) // 发表评论（需登录）
		}
	}
}
