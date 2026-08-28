package dto

// AdminCategoryListRequest 后台分类列表请求
type AdminCategoryListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword"`                                              // 名称模糊搜索
	Status   string `form:"status" binding:"omitempty,oneof=active inactive all"` // 按状态筛选
	MinCount *int   `form:"min_count" binding:"omitempty,min=0"`                  // 最小文章数
	MaxCount *int   `form:"max_count" binding:"omitempty,min=0"`                  // 最大文章数
}

// AdminCategoryCreateRequest 后台创建分类请求
type AdminCategoryCreateRequest struct {
	Name   string `json:"name" binding:"required,max=50"`
	Slug   string `json:"slug" binding:"omitempty,max=60"` // 留空则自动生成
	Desc   string `json:"desc" binding:"omitempty,max=255"`
	Color  string `json:"color" binding:"omitempty,max=10"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"` // 默认 active
}

// AdminCategoryUpdateRequest 后台更新分类请求
type AdminCategoryUpdateRequest struct {
	Name   *string `json:"name" binding:"omitempty,max=50"`
	Slug   *string `json:"slug" binding:"omitempty,max=60"`
	Desc   *string `json:"desc" binding:"omitempty,max=255"`
	Color  *string `json:"color" binding:"omitempty,max=10"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// AdminCategoryStatusRequest 后台修改分类状态请求
type AdminCategoryStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}
