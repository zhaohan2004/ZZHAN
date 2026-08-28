package admin

import (
	"ZZHAN/internal/middleware"
	"ZZHAN/internal/repository"
	"ZZHAN/internal/service"
	"ZZHAN/pkg/response"

	"github.com/gin-gonic/gin"
)

type UploadController struct {
	uploadSvc service.UploadService
	redisRepo repository.RedisRepository
}

func NewUploadController(uploadSvc service.UploadService, redisRepo repository.RedisRepository) *UploadController {
	return &UploadController{uploadSvc: uploadSvc, redisRepo: redisRepo}
}

// UploadImage 单张图片上传
func (c *UploadController) UploadImage(ctx *gin.Context) {
	userID := middleware.GetUserID(ctx)
	if userID == 0 {
		response.Unauthorized(ctx, "未登录")
		return
	}

	file, err := ctx.FormFile("file")
	if err != nil {
		response.BadRequest(ctx, "参数错误: 请选择要上传的图片")
		return
	}

	resp, err := c.uploadSvc.UploadImages(ctx, file)
	if err != nil {
		response.BizError(ctx, err)
		return
	}
	response.Success(ctx, resp)
}
