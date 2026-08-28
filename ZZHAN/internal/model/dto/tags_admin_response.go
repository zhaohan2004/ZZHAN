package dto

// AdminTagItem 后台标签列表项
type AdminTagItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	ArticleCount int64  `json:"count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// AdminTagListResponse 后台标签列表响应
type AdminTagListResponse struct {
	List  []AdminTagItem `json:"list"`
	Total int64          `json:"total"`
}
