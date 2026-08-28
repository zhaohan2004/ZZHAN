package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *AdminArticlesController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/articles")
	admin.Use(middleware.Auth(c.redisRepo))
	{
		admin.GET("", c.List)                    // GET /api/v1/admin/articles
		admin.GET("/:id", c.GetByID)             // GET /api/v1/admin/articles/:id
		admin.POST("", c.Create)                 // POST /api/v1/admin/articles
		admin.PUT("/:id", c.Update)              // PUT /api/v1/admin/articles/:id
		admin.DELETE("/:id", c.Delete)           // DELETE /api/v1/admin/articles/:id
		admin.PUT("/:id/status", c.UpdateStatus) // PUT /api/v1/admin/articles/:id/status
	}
}
