package web

import (
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"github.com/gin-gonic/gin"
)

type StatsController struct {
	statsService service.StatsService
}

func NewStatsController(statsService service.StatsService) *StatsController {
	return &StatsController{statsService: statsService}
}

// GetStats 获取站点统计
func (c *StatsController) GetStats(ctx *gin.Context) {
	resp, err := c.statsService.GetStats(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取站点统计失败")
		return
	}
	response.Success(ctx, resp)
}
