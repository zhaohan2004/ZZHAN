package web

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册标签相关路由
func (c *TagsController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/tags", c.GetTags) // GET /api/v1/tags - 标签列表
}
