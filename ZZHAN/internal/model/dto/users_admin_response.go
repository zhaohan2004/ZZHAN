package dto

// UsersAdminItem 后台用户列表项
type UsersAdminItem struct {
	ID          int64  `json:"id"`
	Provider    string `json:"provider"`      // 登录方式
	Openid      string `json:"openid"`        // 平台用户标识
	Nickname    string `json:"nickname"`      // 昵称
	Avatar      string `json:"avatar"`        // 头像 URL
	Status      int8   `json:"status"`        // 状态：1 正常 / 0 禁用
	LastLoginAt string `json:"last_login_at"` // 最后登录时间
	CreatedAt   string `json:"created_at"`    // 创建时间
}

// UsersAdminListResponse 后台用户列表响应
type UsersAdminListResponse struct {
	List  []UsersAdminItem `json:"list"`
	Total int64            `json:"total"`
}
