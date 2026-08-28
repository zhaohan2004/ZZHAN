package dto

// AdminTagListRequest 后台标签列表请求
type AdminTagListRequest struct {
	Page     int    `form:"page" binding:"omitempty,min=1"`
	PageSize int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword  string `form:"keyword"`                                              // 名称模糊搜索
	Status   string `form:"status" binding:"omitempty,oneof=active inactive all"` // 按状态筛选
	MinCount *int   `form:"min_count" binding:"omitempty,min=0"`                  // 最小使用次数
	MaxCount *int   `form:"max_count" binding:"omitempty,min=0"`                  // 最大使用次数
}

// AdminTagCreateRequest 后台创建标签请求
type AdminTagCreateRequest struct {
	Name   string `json:"name" binding:"required,max=50"`
	Status string `json:"status" binding:"omitempty,oneof=active inactive"` // 默认 active
}

// AdminTagUpdateRequest 后台更新标签请求
type AdminTagUpdateRequest struct {
	Name   *string `json:"name" binding:"omitempty,max=50"`
	Status *string `json:"status" binding:"omitempty,oneof=active inactive"`
}

// AdminTagStatusRequest 后台修改标签状态请求
type AdminTagStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=active inactive"`
}
