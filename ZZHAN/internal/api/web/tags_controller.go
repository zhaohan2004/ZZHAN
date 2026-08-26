package web

import (
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TagsController struct {
	tagsService service.TagsService
}

func NewTagsController(tagsService service.TagsService) *TagsController {
	return &TagsController{tagsService: tagsService}
}

// GetTags 获取标签列表
// GET /api/v1/tags
func (c *TagsController) GetTags(ctx *gin.Context) {
	tags, err := c.tagsService.GetAll(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取标签列表失败")
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    tags,
	})
}
