package dto

// LoginAdminRequest 后台登录请求
type LoginAdminRequest struct {
	UserName  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required,min=6"`
	CaptchaId string `json:"captcha_id" binding:"required"`
	Captcha   string `json:"captcha" binding:"required"`
}

// UpdateAdminProfileRequest 更新管理员资料请求
type UpdateAdminProfileRequest struct {
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Password string `json:"password"` // 可选，为空则不修改密码
}
