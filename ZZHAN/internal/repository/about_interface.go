package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// AboutRepository 关于我仓储接口
type AboutRepository interface {
	// GetAbout 获取关于我数据（技能列表）
	GetAbout(ctx context.Context) (*dto.AboutResponse, error)
}
