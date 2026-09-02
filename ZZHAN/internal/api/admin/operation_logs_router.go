package admin

import (
	"ZZHAN/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (c *AdminOperationLogsController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/operation-logs")
	admin.Use(middleware.Auth(c.redisRepo, "admin"))
	{
		admin.GET("", c.List) // GET /api/v1/admin/operation-logs
	}
}
