package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *AdminUsersController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/users")
	admin.Use(middleware.Auth(c.redisRepo, "admin"))
	admin.Use(middleware.OperationLog(c.db))
	{
		admin.GET("", c.List)                    // GET /api/v1/admin/users
		admin.GET("/:id", c.GetByID)             // GET /api/v1/admin/users/:id
		admin.PUT("/:id/status", c.UpdateStatus) // PUT /api/v1/admin/users/:id/status
		admin.DELETE("/:id", c.Delete)           // DELETE /api/v1/admin/users/:id
	}
}
