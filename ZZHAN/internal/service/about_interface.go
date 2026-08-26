package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// AboutService 关于我业务接口
type AboutService interface {
	// GetAbout 获取关于我数据
	GetAbout(ctx context.Context) (*dto.AboutResponse, error)
}
