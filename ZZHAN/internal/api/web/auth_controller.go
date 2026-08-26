package web

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
	redisRepo   repository.RedisRepository
}

func NewAuthController(authService service.AuthService, redisRepo repository.RedisRepository) *AuthController {
	return &AuthController{authService: authService, redisRepo: redisRepo}
}

// GitHubLogin GitHub OAuth 登录
func (c *AuthController) GitHubLogin(ctx *gin.Context) {
	var req dto.LoginRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"message": "请求参数错误",
		})
		return
	}

	resp, err := c.authService.GitHubLogin(ctx.Request.Context(), req.Code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "登录失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "登录成功",
		"data":    resp,
	})
}

// RefreshToken 刷新 access_token
func (c *AuthController) RefreshToken(ctx *gin.Context) {
	var req dto.RefreshToken
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"message": "请求参数错误",
		})
		return
	}

	resp, err := c.authService.RefreshToken(ctx.Request.Context(), req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    40100,
			"message": "刷新失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "刷新成功",
		"data":    resp,
	})
}

// Logout 退出登录
func (c *AuthController) Logout(ctx *gin.Context) {
	// 从 Header 获取 access_token
	accessToken := ctx.GetHeader("Authorization")
	if accessToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"message": "缺少 access_token",
		})
		return
	}

	// 去掉 "Bearer " 前缀
	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	err := c.authService.Logout(ctx.Request.Context(), accessToken)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "退出失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "退出成功",
	})
}

// GetCurrentUser 获取当前登录用户信息
// GET /auth/me
func (c *AuthController) GetCurrentUser(ctx *gin.Context) {
	// 从上下文获取用户ID（由 Auth 中间件设置）
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    40100,
			"message": "未登录",
		})
		return
	}

	resp, err := c.authService.GetCurrentUser(ctx.Request.Context(), userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "获取用户信息失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    resp,
	})
}
