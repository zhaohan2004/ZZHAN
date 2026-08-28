package dto

// AdminArticleListRequest 后台文章列表请求
type AdminArticleListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword"`                                                   // 搜索关键词（标题）
	Category string `form:"category"`                                                  // 按分类名称筛选
	Tag      string `form:"tag"`                                                       // 按标签名称筛选
	Status   string `form:"status" binding:"omitempty,oneof=published draft down all"` // 按状态筛选
}

// AdminArticleCreateRequest 后台创建文章请求
type AdminArticleCreateRequest struct {
	Title       string   `json:"title" binding:"required,max=200"`
	Summary     string   `json:"summary" binding:"omitempty,max=500"`
	CoverImage  string   `json:"cover_image" binding:"omitempty,max=255"`
	Category    string   `json:"category" binding:"required,max=50"` // 分类名称
	Tags        []string `json:"tags"`                               // 标签名称列表
	Content     string   `json:"content" binding:"required"`
	Status      string   `json:"status" binding:"required,oneof=published draft down"`
	PublishedAt string   `json:"published_at"` // 格式：2006-01-02
}

// AdminArticleUpdateRequest 后台更新文章请求
type AdminArticleUpdateRequest struct {
	Title       string   `json:"title" binding:"required,max=200"`
	Summary     string   `json:"summary" binding:"omitempty,max=500"`
	CoverImage  string   `json:"cover_image" binding:"omitempty,max=255"`
	Category    string   `json:"category" binding:"required,max=50"` // 分类名称
	Tags        []string `json:"tags"`                               // 标签名称列表
	Content     string   `json:"content" binding:"required"`
	Status      string   `json:"status" binding:"required,oneof=published draft down"`
	PublishedAt string   `json:"published_at"` // 格式：2006-01-02
}

// AdminArticleStatusRequest 后台修改文章状态请求
type AdminArticleStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=published draft down"`
}
