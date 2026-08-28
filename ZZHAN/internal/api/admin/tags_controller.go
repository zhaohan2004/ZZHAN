package admin

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminTagsController struct {
	tagsAdminService service.TagsAdminService
	redisRepo        repository.RedisRepository
}

func NewAdminTagsController(tagsAdminService service.TagsAdminService, redisRepo repository.RedisRepository) *AdminTagsController {
	return &AdminTagsController{
		tagsAdminService: tagsAdminService,
		redisRepo:        redisRepo,
	}
}

// List 标签列表
func (c *AdminTagsController) List(ctx *gin.Context) {
	var req dto.AdminTagListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	resp, err := c.tagsAdminService.AdminList(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalError(ctx, "获取标签列表失败")
		return
	}

	response.Success(ctx, resp)
}

// GetByID 获取标签详情
func (c *AdminTagsController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "标签 ID 格式错误")
		return
	}

	item, err := c.tagsAdminService.AdminGetByID(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound(ctx, "标签不存在")
		return
	}

	response.Success(ctx, item)
}

// Create 创建标签
func (c *AdminTagsController) Create(ctx *gin.Context) {
	var req dto.AdminTagCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	item, err := c.tagsAdminService.AdminCreate(ctx.Request.Context(), &req)
	if err != nil {
		if isDuplicateError(err) {
			response.BadRequest(ctx, "标签名称已存在")
			return
		}
		response.InternalError(ctx, "创建标签失败")
		return
	}

	response.Success(ctx, item)
}

// Update 更新标签
func (c *AdminTagsController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "标签 ID 格式错误")
		return
	}

	var req dto.AdminTagUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	item, err := c.tagsAdminService.AdminUpdate(ctx.Request.Context(), id, &req)
	if err != nil {
		if isDuplicateError(err) {
			response.BadRequest(ctx, "标签名称已存在")
			return
		}
		response.InternalError(ctx, "更新标签失败")
		return
	}

	response.Success(ctx, item)
}

// UpdateStatus 修改标签状态
func (c *AdminTagsController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "标签 ID 格式错误")
		return
	}

	var req dto.AdminTagStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	if err := c.tagsAdminService.AdminUpdateStatus(ctx.Request.Context(), id, req.Status); err != nil {
		response.InternalError(ctx, "修改状态失败")
		return
	}

	response.Success(ctx, nil)
}

// Delete 删除标签
func (c *AdminTagsController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "标签 ID 格式错误")
		return
	}

	if err := c.tagsAdminService.AdminDelete(ctx.Request.Context(), id); err != nil {
		response.InternalError(ctx, "删除标签失败")
		return
	}

	response.Success(ctx, nil)
}
