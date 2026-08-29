package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type dashboardAdminService struct {
	dashboardAdminRepo repository.DashboardAdminRepository
}

// NewDashboardAdminService 创建后台仪表盘业务实例
func NewDashboardAdminService(dashboardAdminRepo repository.DashboardAdminRepository) DashboardAdminService {
	return &dashboardAdminService{dashboardAdminRepo: dashboardAdminRepo}
}

// GetStats 获取统计数据
func (s *dashboardAdminService) GetStats(ctx context.Context) (map[string]dto.StatItem, error) {
	return s.dashboardAdminRepo.GetStats(ctx)
}

// GetArticles 获取最新发布文章
func (s *dashboardAdminService) GetArticles(ctx context.Context) ([]dto.RecentPostItem, error) {
	return s.dashboardAdminRepo.GetArticles(ctx)
}

// GetComments 获取最新5条评论
func (s *dashboardAdminService) GetComments(ctx context.Context) ([]dto.DashboardCommentItem, error) {
	return s.dashboardAdminRepo.GetComments(ctx)
}

// GetOperations 获取最新5条操作记录
func (s *dashboardAdminService) GetOperations(ctx context.Context) ([]dto.DashboardOperationItem, error) {
	return s.dashboardAdminRepo.GetOperations(ctx)
}
