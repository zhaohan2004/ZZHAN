package service

import (
	"ZZHAN/internal/model/dto"
	"context"
	"mime/multipart"
)

type UploadService interface {
	// UploadImages 上传图片，返回 URL
	UploadImages(ctx context.Context, files *multipart.FileHeader) (*dto.UploadImageResponse, error)
}
