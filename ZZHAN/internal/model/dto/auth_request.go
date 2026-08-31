package dto

// LoginRequest GitHub响应码
type LoginRequest struct {
	Code        string `json:"code" binding:"required"`
	RedirectURI string `json:"redirect_uri"`
}

// RefreshToken 刷新AccessToken
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
