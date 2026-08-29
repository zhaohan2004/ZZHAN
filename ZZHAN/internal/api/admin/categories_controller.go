package admin

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminCategoriesController struct {
	categoriesAdminService service.CategoriesAdminService
	redisRepo              repository.RedisRepository
	db                     *gorm.DB
}

func NewAdminCategoriesController(categoriesAdminService service.CategoriesAdminService, redisRepo repository.RedisRepository, db *gorm.DB) *AdminCategoriesController {
	return &AdminCategoriesController{
		categoriesAdminService: categoriesAdminService,
		redisRepo:              redisRepo,
		db:                     db,
	}
}

// List 分类列表
func (c *AdminCategoriesController) List(ctx *gin.Context) {
	var req dto.AdminCategoryListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	resp, err := c.categoriesAdminService.AdminList(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalError(ctx, "获取分类列表失败")
		return
	}

	response.Success(ctx, resp)
}

// GetByID 获取分类详情
func (c *AdminCategoriesController) GetByID(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "分类 ID 格式错误")
		return
	}

	item, err := c.categoriesAdminService.AdminGetByID(ctx.Request.Context(), id)
	if err != nil {
		response.NotFound(ctx, "分类不存在")
		return
	}

	response.Success(ctx, item)
}

// Create 创建分类
func (c *AdminCategoriesController) Create(ctx *gin.Context) {
	var req dto.AdminCategoryCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	item, err := c.categoriesAdminService.AdminCreate(ctx.Request.Context(), &req)
	if err != nil {
		if isDuplicateError(err) {
			response.BadRequest(ctx, "分类名称已存在")
			return
		}
		response.InternalError(ctx, "创建分类失败")
		return
	}

	response.Success(ctx, item)
}

// Update 更新分类
func (c *AdminCategoriesController) Update(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "分类 ID 格式错误")
		return
	}

	var req dto.AdminCategoryUpdateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	item, err := c.categoriesAdminService.AdminUpdate(ctx.Request.Context(), id, &req)
	if err != nil {
		if isDuplicateError(err) {
			response.BadRequest(ctx, "分类名称已存在")
			return
		}
		response.InternalError(ctx, "更新分类失败")
		return
	}

	response.Success(ctx, item)
}

// UpdateStatus 修改分类状态
func (c *AdminCategoriesController) UpdateStatus(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "分类 ID 格式错误")
		return
	}

	var req dto.AdminCategoryStatusRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "请求参数错误")
		return
	}

	if err := c.categoriesAdminService.AdminUpdateStatus(ctx.Request.Context(), id, req.Status); err != nil {
		response.InternalError(ctx, "修改状态失败")
		return
	}

	response.Success(ctx, nil)
}

// isDuplicateError 判断是否为唯一约束冲突错误
func isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "Duplicate entry") ||
		strings.Contains(err.Error(), "duplicate key") ||
		strings.Contains(err.Error(), "UNIQUE constraint")
}

// Delete 删除分类
func (c *AdminCategoriesController) Delete(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		response.BadRequest(ctx, "分类 ID 格式错误")
		return
	}

	if err := c.categoriesAdminService.AdminDelete(ctx.Request.Context(), id); err != nil {
		response.InternalError(ctx, "删除分类失败")
		return
	}

	response.Success(ctx, nil)
}
