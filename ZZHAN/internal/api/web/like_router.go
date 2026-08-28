package web

import (
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes 注册点赞相关路由
func (c *LikeController) RegisterRoutes(r *gin.RouterGroup) {
	like := r.Group("/like")
	{
		// 点赞需要登录
		auth := like.Group("")
		auth.Use(middleware.Auth(c.redisRepo))
		{
			auth.POST("/article/:slug", c.ArticleLike) // POST /api/v1/like/article/:slug - 文章点赞
			auth.POST("/comment/:id", c.CommentLike)   // POST /api/v1/like/comment/:id - 评论点赞
		}
	}
}
