package admin

import (
	"net/http"

	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/b64c"

	"github.com/gin-gonic/gin"
)

type AdminAuthController struct {
	adminAuthService service.AdminAuthService
	redisRepo        repository.RedisRepository
}

func NewAdminAuthController(adminAuthService service.AdminAuthService, redisRepo repository.RedisRepository) *AdminAuthController {
	return &AdminAuthController{adminAuthService: adminAuthService, redisRepo: redisRepo}
}

// GetCaptcha 生成验证码
func (c *AdminAuthController) GetCaptcha(ctx *gin.Context) {
	id, captchaImage, err := b64c.Generate()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "生成验证码失败",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data": dto.CaptchaResponse{
			CaptchaId:    id,
			CaptchaImage: captchaImage,
		},
	})
}

// Login 后台管理员登录
func (c *AdminAuthController) Login(ctx *gin.Context) {
	var req dto.LoginAdminRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"message": "请求参数错误",
		})
		return
	}

	// 校验验证码
	if err := b64c.Identify(req.CaptchaId, req.Captcha); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    3004,
			"message": "验证码错误",
		})
		return
	}

	// 登录
	resp, err := c.adminAuthService.Login(ctx.Request.Context(), &req)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    40100,
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

// Logout 退出登录
func (c *AdminAuthController) Logout(ctx *gin.Context) {
	accessToken := ctx.GetHeader("Authorization")
	if accessToken == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"message": "缺少 access_token",
		})
		return
	}

	if len(accessToken) > 7 && accessToken[:7] == "Bearer " {
		accessToken = accessToken[7:]
	}

	if err := c.adminAuthService.Logout(ctx.Request.Context(), accessToken); err != nil {
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

// GetProfile 获取管理员资料
func (c *AdminAuthController) GetProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    40100,
			"message": "未登录",
		})
		return
	}

	resp, err := c.adminAuthService.GetProfile(ctx.Request.Context(), userID.(int))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "获取资料失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "成功",
		"data":    resp,
	})
}

// UpdateProfile 更新管理员资料
func (c *AdminAuthController) UpdateProfile(ctx *gin.Context) {
	userID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"code":    40100,
			"message": "未登录",
		})
		return
	}

	var req dto.UpdateAdminProfileRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"code":    40000,
			"message": "请求参数错误",
		})
		return
	}

	resp, err := c.adminAuthService.UpdateProfile(ctx.Request.Context(), userID.(int), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"code":    50000,
			"message": "更新失败：" + err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"code":    0,
		"message": "更新成功",
		"data":    resp,
	})
}
