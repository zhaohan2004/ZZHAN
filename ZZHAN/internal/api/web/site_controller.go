package web

import (
	"ZZHAN/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SiteController struct {
	siteService service.SiteService
}

func NewSiteController(siteService service.SiteService) *SiteController {
	return &SiteController{siteService: siteService}
}

// GetSite 获取站点信息（首页加关于）
func (c *SiteController) GetSite(ctx *gin.Context) {
	resp, err := c.siteService.GetSite(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "获取站点信息失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    resp,
	})
}
