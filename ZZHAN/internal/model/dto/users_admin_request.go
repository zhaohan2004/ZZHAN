package dto

// UsersAdminListRequest 后台用户列表请求
type UsersAdminListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword   string `form:"keyword"`                                              // 搜索关键词（昵称）
	Status    string `form:"status" binding:"omitempty,oneof=active inactive all"` // 按状态筛选
	StartDate string `form:"start_date"`                                           // 开始日期（2006-01-02）
	EndDate   string `form:"end_date"`                                             // 结束日期（2006-01-02）
}

// UsersAdminUpdateStatusRequest 后台修改用户状态请求
type UsersAdminUpdateStatusRequest struct {
	Status int8 `json:"status" binding:"required,oneof=0 1"` // 0 禁用 / 1 正常
}
