package admin

import (
	"ZZHAN/internal/middleware"
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AdminArticlesController struct {
	articlesAdminService service.ArticlesAdminService
	redisRepo            repository.RedisRepository
}

func NewAdminArticlesController(articlesAdminService service.ArticlesAdminService, redisRepo repository.RedisRepository) *AdminArticlesController {
	return &AdminArticlesController{
		articlesAdminService: articlesAdminService,
		redisRepo:            redisRepo,
	}
}

// List 文章列表
func (c *AdminArticlesController) List(ctx *gin.Context) {
	var req dto.AdminArticleListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	resp, err := c.articlesAdminService.AdminList(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalError(ctx, "获取文章列表失败")
		return
	}

	response.Success(ctx, resp)
}

// GetByID 获取文章详情
func (c *AdminArticlesController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "文章 ID 格式错误")
		return
	}

	detail, err := c.articlesAdminService.AdminGetByID(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound(ctx, "文章不存在")
		return
	}

	response.Success(ctx, detail)
}

// Create 创建文章
func (c *AdminArticlesController) Create(ctx *gin.Context) {
	var req dto.AdminArticleCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	authorID := int64(middleware.GetUserID(ctx))
	if authorID == 0 {
		response.Unauthorized(ctx, "未登录")
		return
	}

	detail, err := c.articlesAdminService.AdminCreate(ctx.Request.Context(), &req, authorID)
	if err != nil {
		response.InternalError(ctx, "创建文章失败："+err.Error())
		return
	}

	response.Success(ctx, detail)
}

// Update 更新文章
func (c *AdminArticlesController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "文章 ID 格式错误")
		return
	}

	var req dto.AdminArticleUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	detail, err := c.articlesAdminService.AdminUpdate(ctx.Request.Context(), id, &req)
	if err != nil {
		response.InternalError(ctx, "更新文章失败："+err.Error())
		return
	}

	response.Success(ctx, detail)
}

// Delete 删除文章
func (c *AdminArticlesController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "文章 ID 格式错误")
		return
	}

	if err := c.articlesAdminService.AdminDelete(ctx.Request.Context(), id); err != nil {
		response.InternalError(ctx, "删除文章失败")
		return
	}

	response.Success(ctx, nil)
}

// UpdateStatus 修改文章状态
func (c *AdminArticlesController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "文章 ID 格式错误")
		return
	}

	var req dto.AdminArticleStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	detail, err := c.articlesAdminService.AdminUpdateStatus(ctx.Request.Context(), id, req.Status)
	if err != nil {
		response.BadRequest(ctx, err.Error())
		return
	}

	response.Success(ctx, detail)
}
