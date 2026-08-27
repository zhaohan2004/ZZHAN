package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type statsService struct {
	statsRepo repository.StatsRepository
}

// NewStatsService 创建统计业务实例
func NewStatsService(statsRepo repository.StatsRepository) StatsService {
	return &statsService{statsRepo: statsRepo}
}

// GetStats 获取站点统计数据
func (s *statsService) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	return s.statsRepo.GetStats(ctx)
}
