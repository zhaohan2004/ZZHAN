package admin

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminCommentsController struct {
	commentsAdminService service.CommentsAdminService
	redisRepo            repository.RedisRepository
	db                   *gorm.DB
}

func NewAdminCommentsController(commentsAdminService service.CommentsAdminService, redisRepo repository.RedisRepository, db *gorm.DB) *AdminCommentsController {
	return &AdminCommentsController{
		commentsAdminService: commentsAdminService,
		redisRepo:            redisRepo,
		db:                   db,
	}
}

// List 评论列表
func (c *AdminCommentsController) List(ctx *gin.Context) {
	var req dto.CommentsAdminListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	resp, err := c.commentsAdminService.AdminList(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalError(ctx, "获取评论列表失败")
		return
	}

	response.Success(ctx, resp)
}

// GetByID 获取评论详情
func (c *AdminCommentsController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "评论 ID 格式错误")
		return
	}

	detail, err := c.commentsAdminService.AdminGetByID(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound(ctx, "评论不存在")
		return
	}

	response.Success(ctx, detail)
}

// UpdateStatus 修改评论状态
func (c *AdminCommentsController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "评论 ID 格式错误")
		return
	}

	var req dto.CommentsAdminUpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	detail, err := c.commentsAdminService.AdminUpdateStatus(ctx.Request.Context(), id, req.Status)
	if err != nil {
		response.InternalError(ctx, "修改评论状态失败")
		return
	}

	response.Success(ctx, detail)
}

// Delete 删除评论
func (c *AdminCommentsController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "评论 ID 格式错误")
		return
	}

	if err := c.commentsAdminService.AdminDelete(ctx.Request.Context(), id); err != nil {
		response.InternalError(ctx, "删除评论失败")
		return
	}

	response.Success(ctx, nil)
}
