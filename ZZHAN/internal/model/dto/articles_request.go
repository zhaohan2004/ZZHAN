package dto

import "ZZHAN/pkg/response"

// ArticleListRequest 文章列表请求
// 嵌入 response.PageRequest 复用分页参数
type ArticleListRequest struct {
	response.PageRequest        // 嵌入分页参数：Page, Size
	CategoryID           int64  `form:"category_id"` // 按分类筛选（可选）
	TagID                int64  `form:"tag_id"`      // 按标签筛选（可选）
	Keyword              string `form:"keyword"`     // 搜索关键词（可选）
}

// ArticleDetailRequest 文章详情请求（通过 slug 获取）
type ArticleDetailRequest struct {
	Slug string `uri:"slug" binding:"required"` // URL 中的文章别名
}
