package dto

// OperationLogsAdminItem 后台操作日志列表项
type OperationLogsAdminItem struct {
	ID        int64  `json:"id"`
	AdminID   *int64 `json:"admin_id"`
	AdminName string `json:"admin_name"` // 管理员昵称
	Action    string `json:"action"`
	Target    string `json:"target"`
	CreatedAt string `json:"created_at"`
}

// OperationLogsAdminListResponse 后台操作日志列表响应
type OperationLogsAdminListResponse struct {
	List  []OperationLogsAdminItem `json:"list"`
	Total int64                    `json:"total"`
}
