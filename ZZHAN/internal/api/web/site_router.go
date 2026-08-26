package web

import (
	"github.com/gin-gonic/gin"
)

func (c *SiteController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/site", c.GetSite)
}
