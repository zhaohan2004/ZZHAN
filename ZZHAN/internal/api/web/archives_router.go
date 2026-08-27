package web

import "github.com/gin-gonic/gin"

func (c *ArchivesController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/archives", c.GetArchives)
}
