package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// StatsRepository 统计仓储接口
type StatsRepository interface {
	// GetStats 获取站点统计数据
	GetStats(ctx context.Context) (*dto.StatsResponse, error)
}
