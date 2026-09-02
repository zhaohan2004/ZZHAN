package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *SiteAdminController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin/settings")
	admin.Use(middleware.Auth(c.redisRepo, "admin"))
	admin.Use(middleware.OperationLog(c.db))
	{
		admin.GET("", c.GetSettings)
		admin.PUT("", c.UpdateSettings)
	}
}
