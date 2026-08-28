package dto

// AdminArticleItem 后台文章列表项
type AdminArticleItem struct {
	ID         int64    `json:"id"`
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	CoverImage string   `json:"cover_image"`
	Category   string   `json:"category"` // 分类名称
	Tags       []string `json:"tags"`     // 标签名称列表
	Status     string   `json:"status"`
	Views      int32    `json:"views"`
	Date       string   `json:"date"`    // 发布时间
	Updated    string   `json:"updated"` // 更新时间
}

// AdminArticleDetail 后台文章详情
type AdminArticleDetail struct {
	ID         int64    `json:"id"`
	Slug       string   `json:"slug"`
	Title      string   `json:"title"`
	Summary    string   `json:"summary"`
	CoverImage string   `json:"cover_image"`
	Category   string   `json:"category"` // 分类名称
	Tags       []string `json:"tags"`     // 标签名称列表
	Status     string   `json:"status"`
	Views      int32    `json:"views"`
	Date       string   `json:"date"`    // 发布时间
	Updated    string   `json:"updated"` // 更新时间
	Content    string   `json:"content"` // Markdown 正文
}

// AdminArticleListResponse 后台文章列表响应
type AdminArticleListResponse struct {
	List  []AdminArticleItem `json:"list"`
	Total int64              `json:"total"`
}
