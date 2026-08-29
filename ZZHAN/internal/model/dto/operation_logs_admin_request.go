package dto

// OperationLogsAdminListRequest 后台操作日志列表请求
type OperationLogsAdminListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Action    string `form:"action"`     // 按操作类型筛选（新建/更新/删除/发布/存为草稿/下架/启用/禁用/解封/封禁）
	Target    string `form:"target"`     // 按操作对象模糊搜索
	StartDate string `form:"start_date"` // 开始日期（2006-01-02）
	EndDate   string `form:"end_date"`   // 结束日期（2006-01-02）
}
