package web

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"

	"github.com/gin-gonic/gin"
)

type ArticlesController struct {
	articlesService service.ArticlesService
}

func NewArticlesController(articlesService service.ArticlesService) *ArticlesController {
	return &ArticlesController{articlesService: articlesService}
}

// GetArticleList 获取文章列表
// GET /api/v1/articles?page=1&size=10&category_id=1&tag_id=1&keyword=xxx
func (c *ArticlesController) GetArticleList(ctx *gin.Context) {
	var req dto.ArticleListRequest

	// 绑定查询参数
	if err := ctx.ShouldBindQuery(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	// 调用 service
	result, err := c.articlesService.GetPublishedList(ctx.Request.Context(), &req)
	if err != nil {
		response.InternalError(ctx, "获取文章列表失败")
		return
	}

	response.Success(ctx, result)
}

// GetArticleDetail 获取文章详情
// GET /api/v1/articles/:slug
func (c *ArticlesController) GetArticleDetail(ctx *gin.Context) {
	var req dto.ArticleDetailRequest

	// 绑定路径参数
	if err := ctx.ShouldBindUri(&req); err != nil {
		response.BadRequest(ctx, "参数错误: "+err.Error())
		return
	}

	// 获取客户端 IP（用于浏览量去重）
	clientIP := ctx.ClientIP()

	// 调用 service
	article, err := c.articlesService.GetBySlug(ctx.Request.Context(), req.Slug, clientIP)
	if err != nil {
		// 判断是否是记录不存在
		if err.Error() == "record not found" {
			response.NotFound(ctx, "文章不存在")
			return
		}
		response.InternalError(ctx, "获取文章详情失败")
		return
	}

	response.Success(ctx, article)
}
