package web

import (
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"github.com/gin-gonic/gin"
)

type ArchivesController struct {
	archivesService service.ArchivesService
}

func NewArchivesController(archivesService service.ArchivesService) *ArchivesController {
	return &ArchivesController{archivesService: archivesService}
}

// GetArchives 获取归档列表
func (c *ArchivesController) GetArchives(ctx *gin.Context) {
	resp, err := c.archivesService.GetArchives(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取归档列表失败")
		return
	}
	response.Success(ctx, resp)
}
