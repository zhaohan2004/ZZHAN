package web

import (
	"ZZHAN/internal/middleware"
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CommentsController struct {
	commentsService service.CommentsService
	redisRepo       repository.RedisRepository
	commentRepo     repository.CommentsRepository
}

func NewCommentsController(commentsService service.CommentsService, redisRepo repository.RedisRepository, commentRepo repository.CommentsRepository) *CommentsController {
	return &CommentsController{
		commentsService: commentsService,
		redisRepo:       redisRepo,
		commentRepo:     commentRepo,
	}
}

// GetComments 获取文章评论列表
// GET /api/v1/comments/:slug?page=1&page_size=10
func (c *CommentsController) GetComments(ctx *gin.Context) {
	slug := ctx.Param("slug")

	var req dto.CommentListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := c.commentsService.GetByArticleSlug(ctx.Request.Context(), slug, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(ctx, "获取评论列表失败")
		return
	}

	response.Success(ctx, result)
}

// GetReplies 获取评论的回复列表
// GET /api/v1/comments/:id/replies?page=1&page_size=10
func (c *CommentsController) GetReplies(ctx *gin.Context) {
	idStr := ctx.Param("id")
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "评论ID格式错误")
		return
	}

	var req dto.CommentListRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := c.commentsService.GetReplies(ctx.Request.Context(), commentID, req.Page, req.PageSize)
	if err != nil {
		response.InternalError(ctx, "获取回复列表失败")
		return
	}

	response.Success(ctx, result)
}

// PostComment 发表评论
// POST /api/v1/comments/:slug
func (c *CommentsController) PostComment(ctx *gin.Context) {
	slug := ctx.Param("slug")

	var req dto.CommentCreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	// 从上下文获取用户信息（由 Auth 中间件写入）
	userID := int64(middleware.GetUserID(ctx))
	if userID == 0 {
		response.Unauthorized(ctx, "请先登录")
		return
	}

	userName := middleware.GetUsername(ctx)

	userAvatar, err := c.commentRepo.GetUserAvatarByID(ctx, userID)
	if err != nil {
		response.BadRequest(ctx, "获取头像失败")
		return
	}

	result, err := c.commentsService.Create(ctx.Request.Context(), slug, userID, userName, userAvatar, ctx.ClientIP(), &req)
	if err != nil {
		response.InternalError(ctx, "发表评论失败: "+err.Error())
		return
	}

	response.Success(ctx, result)
}
