package web

import (
	"ZZHAN/internal/middleware"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type LikeController struct {
	likeService service.LikeService
	redisRepo   repository.RedisRepository
}

func NewLikeController(likeService service.LikeService, redisRepo repository.RedisRepository) *LikeController {
	return &LikeController{
		likeService: likeService,
		redisRepo:   redisRepo,
	}
}

// ArticleLike 文章点赞
func (c *LikeController) ArticleLike(ctx *gin.Context) {
	slug := ctx.Param("slug")

	userID := int64(middleware.GetUserID(ctx))
	if userID == 0 {
		response.Unauthorized(ctx, "请先登录")
		return
	}

	result, err := c.likeService.ArticleLike(ctx.Request.Context(), slug, userID)
	if err != nil {
		response.InternalError(ctx, "点赞失败")
		return
	}

	response.Success(ctx, result)
}

// CommentLike 评论点赞
func (c *LikeController) CommentLike(ctx *gin.Context) {
	idStr := ctx.Param("id")
	commentID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		response.BadRequest(ctx, "评论ID格式错误")
		return
	}

	userID := int64(middleware.GetUserID(ctx))
	if userID == 0 {
		response.Unauthorized(ctx, "请先登录")
		return
	}

	result, err := c.likeService.CommentLike(ctx.Request.Context(), commentID, userID)
	if err != nil {
		response.InternalError(ctx, "点赞失败")
		return
	}

	response.Success(ctx, result)
}
