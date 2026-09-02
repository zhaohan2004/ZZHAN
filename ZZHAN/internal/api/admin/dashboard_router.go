package admin

import (
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (c *AdminDashboardController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/dashboard")
	admin.Use(middleware.Auth(c.redisRepo, "admin"))
	{
		admin.GET("/stats", c.GetStats)
		admin.GET("/articles", c.GetArticles)
		admin.GET("/comments", c.GetComments)
		admin.GET("/operations", c.GetOperations)
	}
}
