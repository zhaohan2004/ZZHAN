package admin

import (
	"ZZHAN/internal/middleware"
	"github.com/gin-gonic/gin"
)

func (c *UploadController) RegisterRoutes(r *gin.RouterGroup) {
	up := r.Group("/upload")
	up.Use(middleware.Auth(c.redisRepo))
	{
		up.POST("/image", c.UploadImage)
	}
}
