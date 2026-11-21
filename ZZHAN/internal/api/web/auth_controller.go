package web

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/config"
	appErrors "ZZHAN/pkg/errors"
	"ZZHAN/pkg/response"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	authService service.AuthService
	userAuth    gin.HandlerFunc
}

func NewAuthController(authService service.AuthService, userAuth gin.HandlerFunc) *AuthController {
	return &AuthController{authService: authService, userAuth: userAuth}
}

func respondAccountError(ctx *gin.Context, err error) bool {
	switch {
	case errors.Is(err, appErrors.ErrUserDisabled):
		response.Error(ctx, http.StatusForbidden, appErrors.CodeUserDisabled, appErrors.ErrUserDisabled.Message)
	case errors.Is(err, appErrors.ErrInternalError):
		_ = ctx.Error(err)
		response.Error(ctx, http.StatusServiceUnavailable, appErrors.CodeInternalError, "暂时无法验证登录状态")
	case errors.Is(err, appErrors.ErrInvalidToken):
		response.Unauthorized(ctx, "登录已失效，请重新登录")
	default:
		return false
	}
	return true
}

// GitHubConfig returns the public OAuth client ID for the authorization request.
func (c *AuthController) GitHubConfig(ctx *gin.Context) {
	cfg := config.Get().GitHub
	clientID := strings.TrimSpace(cfg.ClientID)
	if clientID == "" || strings.TrimSpace(cfg.ClientSecret) == "" {
		response.Error(ctx, http.StatusServiceUnavailable, 50300, "GitHub 登录尚未配置")
		return
	}
	response.Success(ctx, gin.H{"client_id": clientID})
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

	resp, err := c.authService.GitHubLogin(ctx.Request.Context(), req.Code, req.RedirectURI)
	if err != nil {
		if respondAccountError(ctx, err) {
			return
		}
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
		if respondAccountError(ctx, err) {
			return
		}
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
		if respondAccountError(ctx, err) {
			return
		}
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
