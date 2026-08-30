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

type AdminUsersController struct {
	usersAdminService service.UsersAdminService
	redisRepo         repository.RedisRepository
	db                *gorm.DB
}

func NewAdminUsersController(usersAdminService service.UsersAdminService, redisRepo repository.RedisRepository, db *gorm.DB) *AdminUsersController {
	return &AdminUsersController{
		usersAdminService: usersAdminService,
		redisRepo:         redisRepo,
		db:                db,
	}
}

// List 用户列表
func (c *AdminUsersController) List(ctx *gin.Context) {
	var req dto.UsersAdminListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	resp, err := c.usersAdminService.AdminList(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalError(ctx, "获取用户列表失败")
		return
	}

	response.Success(ctx, resp)
}

// GetByID 获取用户详情
func (c *AdminUsersController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户 ID 格式错误")
		return
	}

	detail, err := c.usersAdminService.AdminGetByID(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound(ctx, "用户不存在")
		return
	}

	response.Success(ctx, detail)
}

// UpdateStatus 修改用户状态
func (c *AdminUsersController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户 ID 格式错误")
		return
	}

	var req dto.UsersAdminUpdateStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	detail, err := c.usersAdminService.AdminUpdateStatus(ctx.Request.Context(), id, req.Status)
	if err != nil {
		response.InternalError(ctx, "修改用户状态失败")
		return
	}

	response.Success(ctx, detail)
}

// Delete 删除用户
func (c *AdminUsersController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "用户 ID 格式错误")
		return
	}

	if err := c.usersAdminService.AdminDelete(ctx.Request.Context(), id); err != nil {
		response.InternalError(ctx, "删除用户失败")
		return
	}

	response.Success(ctx, nil)
}
