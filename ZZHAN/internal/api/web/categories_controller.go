package web

import (
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoriesController struct {
	categoriesService service.CategoriesService
}

func NewCategoriesController(categoriesService service.CategoriesService) *CategoriesController {
	return &CategoriesController{categoriesService: categoriesService}
}

// GetCategories 获取分类列表
// GET /api/v1/categories
func (c *CategoriesController) GetCategories(ctx *gin.Context) {
	categories, err := c.categoriesService.GetAll(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取分类列表失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    categories,
	})
}
