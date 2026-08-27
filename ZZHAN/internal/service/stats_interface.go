package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// StatsService 统计业务接口
type StatsService interface {
	// GetStats 获取站点统计数据
	GetStats(ctx context.Context) (*dto.StatsResponse, error)
}
