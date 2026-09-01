package dto

// CaptchaResponse 验证码响应
type CaptchaResponse struct {
	CaptchaId    string `json:"captcha_id"`
	CaptchaImage string `json:"captcha_image"`
}

// AdminProfileResponse 管理员资料响应
type AdminProfileResponse struct {
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}
