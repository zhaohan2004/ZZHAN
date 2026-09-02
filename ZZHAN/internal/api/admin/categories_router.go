package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *AdminCategoriesController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/categories")
	admin.Use(middleware.Auth(c.redisRepo, "admin"))
	admin.Use(middleware.OperationLog(c.db))
	{
		admin.GET("", c.List)                    // GET /api/v1/admin/categories
		admin.GET("/:id", c.GetByID)             // GET /api/v1/admin/categories/:id
		admin.POST("", c.Create)                 // POST /api/v1/admin/categories
		admin.PUT("/:id", c.Update)              // PUT /api/v1/admin/categories/:id
		admin.PUT("/:id/status", c.UpdateStatus) // PUT /api/v1/admin/categories/:id/status
		admin.DELETE("/:id", c.Delete)           // DELETE /api/v1/admin/categories/:id
	}
}
