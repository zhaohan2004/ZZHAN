package web

import "github.com/gin-gonic/gin"

func (c *StatsController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/stats", c.GetStats)
}
