package admin

import (
	"ZZHAN/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type SiteAdminController struct {
	siteService service.SiteService
}

func NewSiteAdminController(siteService service.SiteService) *SiteAdminController {
	return &SiteAdminController{siteService: siteService}
}

// GetSettings 获取所有设置
func (c *SiteAdminController) GetSettings(ctx *gin.Context) {
	settings, err := c.siteService.GetSettings(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "获取设置失败",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    settings,
	})
}

// UpdateSettings 更新设置
func (c *SiteAdminController) UpdateSettings(ctx *gin.Context) {
	var body map[string]string
	if err := ctx.ShouldBindJSON(&body); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    40001,
			"message": "请求参数错误",
		})
		return
	}
	if err := c.siteService.UpdateSettings(ctx.Request.Context(), body); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "保存设置失败",
		})
		return
	}

	// 返回更新后的完整设置
	setting, err := c.siteService.GetSettings(ctx.Request.Context())
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"code":    0,
			"message": "保存成功",
			"data":    body,
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    setting,
	})
}
