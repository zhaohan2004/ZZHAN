package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// DashboardAdminRepository 后台仪表盘仓储接口
type DashboardAdminRepository interface {
	// GetStats 获取统计数据（文章/分类/标签/评论数量）
	GetStats(ctx context.Context) (map[string]dto.StatItem, error)

	// GetArticles 获取最新发布文章
	GetArticles(ctx context.Context) ([]dto.RecentPostItem, error)

	// GetComments 获取最新5条评论
	GetComments(ctx context.Context) ([]dto.DashboardCommentItem, error)

	// GetOperations 获取最新8条操作记录
	GetOperations(ctx context.Context) ([]dto.DashboardOperationItem, error)
}
