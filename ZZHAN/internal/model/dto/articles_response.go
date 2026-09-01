package dto

import "time"

// ArticleListItem 文章列表项（不含正文内容，节省带宽）
type ArticleListItem struct {
	ID           int64     `json:"id"`            // 文章ID
	Title        string    `json:"title"`         // 标题
	Slug         string    `json:"slug"`          // SEO 友好 URL
	Summary      string    `json:"summary"`       // 摘要
	CoverImage   string    `json:"cover_image"`   // 封面图
	CategoryID   int64     `json:"category_id"`   // 分类ID
	CategoryName string    `json:"category_name"` // 分类名称（关联查询）
	AuthorName   string    `json:"author_name"`   // 作者名称
	Tags         []TagItem `json:"tags"`          // 标签列表（关联查询）
	Views        int32     `json:"views"`         // 浏览量
	Likes        int32     `json:"likes"`         // 点赞数
	CommentCount int32     `json:"comment_count"` // 评论数
	PublishedAt  time.Time `json:"published_at"`  // 发布时间
}

// ArticleDetail 文章详情（包含完整正文）
type ArticleDetail struct {
	ID           int64     `json:"id"`            // 文章ID
	Title        string    `json:"title"`         // 标题
	Slug         string    `json:"slug"`          // SEO 友好 URL
	Summary      string    `json:"summary"`       // 摘要
	CoverImage   string    `json:"cover_image"`   // 封面图
	CategoryID   int64     `json:"category_id"`   // 分类ID
	CategoryName string    `json:"category_name"` // 分类名称
	AuthorName   string    `json:"author_name"`   // 作者名称
	Tags         []TagItem `json:"tags"`          // 标签列表
	Content      string    `json:"content"`       // Markdown 正文
	Views        int32     `json:"views"`         // 浏览量
	Likes        int32     `json:"likes"`         // 点赞数
	Liked        bool      `json:"liked"`         // 当前用户是否已点赞
	CommentCount int32     `json:"comment_count"` // 评论数
	PublishedAt  time.Time `json:"published_at"`  // 发布时间
}

// TagItem 标单项（用于文章关联的标签）
type TagItem struct {
	ID   int64  `json:"id"`   // 标签ID
	Name string `json:"name"` // 标签名称
	Slug string `json:"slug"` // 标签 URL 别名
}
