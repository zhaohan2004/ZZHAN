package web

import "github.com/gin-gonic/gin"

// RegisterRoutes 注册分类相关路由
func (c *CategoriesController) RegisterRoutes(r *gin.RouterGroup) {
	r.GET("/categories", c.GetCategories) // GET /api/v1/categories - 分类列表
}
