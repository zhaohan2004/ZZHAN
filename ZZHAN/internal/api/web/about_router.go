package web

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册关于我相关路由
func (c *AboutController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/about", c.GetAbout) // GET /api/v1/about - 关于我
}
