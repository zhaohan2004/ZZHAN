package dto

// AdminCategoryItem 后台分类列表项
type AdminCategoryItem struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Slug         string `json:"slug"`
	Icon         string `json:"icon"`
	Description  string `json:"desc"`
	Color        string `json:"color"`
	Status       string `json:"status"`
	ArticleCount int64  `json:"count"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// AdminCategoryListResponse 后台分类列表响应
type AdminCategoryListResponse struct {
	List  []AdminCategoryItem `json:"list"`
	Total int64               `json:"total"`
}
