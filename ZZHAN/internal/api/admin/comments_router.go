package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *AdminCommentsController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/comments")
	admin.Use(middleware.Auth(c.redisRepo, "admin"))
	admin.Use(middleware.OperationLog(c.db))
	{
		admin.GET("", c.List)                    // GET /api/v1/admin/comments
		admin.GET("/:id", c.GetByID)             // GET /api/v1/admin/comments/:id
		admin.PUT("/:id/status", c.UpdateStatus) // PUT /api/v1/admin/comments/:id/status
		admin.DELETE("/:id", c.Delete)           // DELETE /api/v1/admin/comments/:id
	}
}
