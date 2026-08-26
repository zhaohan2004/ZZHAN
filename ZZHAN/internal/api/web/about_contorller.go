package web

import (
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"

	"github.com/gin-gonic/gin"
)

type AboutController struct {
	aboutService service.AboutService
}

func NewAboutController(aboutService service.AboutService) *AboutController {
	return &AboutController{aboutService: aboutService}
}

// GetAbout 获取关于我数据
// GET /api/v1/about
func (c *AboutController) GetAbout(ctx *gin.Context) {
	resp, err := c.aboutService.GetAbout(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取关于我数据失败")
		return
	}
	response.Success(ctx, resp)
}
