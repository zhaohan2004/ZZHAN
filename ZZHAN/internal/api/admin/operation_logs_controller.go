package admin

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"

	"github.com/gin-gonic/gin"
)

type AdminOperationLogsController struct {
	operationLogsAdminService service.OperationLogsAdminService
	redisRepo                 repository.RedisRepository
}

func NewAdminOperationLogsController(operationLogsAdminService service.OperationLogsAdminService, redisRepo repository.RedisRepository) *AdminOperationLogsController {
	return &AdminOperationLogsController{
		operationLogsAdminService: operationLogsAdminService,
		redisRepo:                 redisRepo,
	}
}

// List 操作日志列表
func (c *AdminOperationLogsController) List(ctx *gin.Context) {
	var req dto.OperationLogsAdminListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	resp, err := c.operationLogsAdminService.AdminList(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalError(ctx, "获取操作日志列表失败")
		return
	}

	response.Success(ctx, resp)
}
