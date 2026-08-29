package admin

import (
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminDashboardController struct {
	dashboardAdminService service.DashboardAdminService
	redisRepo             repository.RedisRepository
}

func NewAdminDashboardController(dashboardAdminService service.DashboardAdminService, redisRepo repository.RedisRepository) *AdminDashboardController {
	return &AdminDashboardController{
		dashboardAdminService: dashboardAdminService,
		redisRepo:             redisRepo,
	}
}

// GetStats 获取统计数据
func (c *AdminDashboardController) GetStats(ctx *gin.Context) {
	data, err := c.dashboardAdminService.GetStats(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取统计数据失败")
		return
	}
	response.Success(ctx, data)
}

// GetArticles 获取最新发布文章
func (c *AdminDashboardController) GetArticles(ctx *gin.Context) {
	data, err := c.dashboardAdminService.GetArticles(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取文章数据失败")
		return
	}
	response.Success(ctx, gin.H{"recent_posts": data})
}

// GetComments 获取最新评论
func (c *AdminDashboardController) GetComments(ctx *gin.Context) {
	data, err := c.dashboardAdminService.GetComments(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取评论数据失败")
		return
	}
	response.Success(ctx, data)
}

// GetOperations 获取最新操作记录
func (c *AdminDashboardController) GetOperations(ctx *gin.Context) {
	data, err := c.dashboardAdminService.GetOperations(ctx.Request.Context())
	if err != nil {
		response.InternalError(ctx, "获取操作记录失败")
		return
	}
	response.Success(ctx, data)
}
