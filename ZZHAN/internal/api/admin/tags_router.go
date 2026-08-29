package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *AdminTagsController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/tags")
	admin.Use(middleware.Auth(c.redisRepo))
	admin.Use(middleware.OperationLog(c.db))
	{
		admin.GET("", c.List)                    // GET /api/v1/admin/tags
		admin.GET("/:id", c.GetByID)             // GET /api/v1/admin/tags/:id
		admin.POST("", c.Create)                 // POST /api/v1/admin/tags
		admin.PUT("/:id", c.Update)              // PUT /api/v1/admin/tags/:id
		admin.PUT("/:id/status", c.UpdateStatus) // PUT /api/v1/admin/tags/:id/status
		admin.DELETE("/:id", c.Delete)           // DELETE /api/v1/admin/tags/:id
	}
}
