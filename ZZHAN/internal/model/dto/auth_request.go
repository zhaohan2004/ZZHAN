package dto

// LoginRequest GitHub响应码
type LoginRequest struct {
	Code string `json:"code" binding:"required"`
}

// RefreshToken 刷新AccessToken
type RefreshToken struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
