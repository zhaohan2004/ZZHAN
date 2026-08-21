package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *SiteAdminController) RegisterRoutes(r *gin.RouterGroup) {
	admin := r.Group("/admin")
	admin.Use(middleware.Auth())
	{
		admin.GET("/settings", c.GetSettings)
		admin.PUT("/settings", c.UpdateSettings)
	}
}
