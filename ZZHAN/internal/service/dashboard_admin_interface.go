package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// DashboardAdminService 后台仪表盘业务接口
type DashboardAdminService interface {
	// GetStats 获取统计数据
	GetStats(ctx context.Context) (map[string]dto.StatItem, error)

	// GetArticles 获取最新发布文章
	GetArticles(ctx context.Context) ([]dto.RecentPostItem, error)

	// GetComments 获取最新5条评论
	GetComments(ctx context.Context) ([]dto.DashboardCommentItem, error)

	// GetOperations 获取最新8条操作记录
	GetOperations(ctx context.Context) ([]dto.DashboardOperationItem, error)
}
