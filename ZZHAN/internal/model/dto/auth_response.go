package dto

// AuthUser 认证用户信息
type AuthUser struct {
	ID       int    `json:"id"`
	Provider string `json:"provider"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string   `json:"access_token"`
	RefreshToken string   `json:"refresh_token"`
	User         AuthUser `json:"user"`
}

// RefreshResponse 刷新token响应
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
