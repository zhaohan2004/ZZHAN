package repository

import (
	"context"

	"ZZHAN/internal/model/dto"
)

// SiteRepository 站点设置仓储接口
type SiteRepository interface {
	// GetSiteResponse 获取站点信息（组装为 SiteResponse）
	GetSiteResponse(ctx context.Context) (*dto.SiteResponse, error)

	// GetSettings 获取所有设置项（key -> value map）
	GetSettings(ctx context.Context) (map[string]string, error)

	// UpdateSettings 批量更新设置项（upsert：存在则更新，不存在则插入）
	UpdateSettings(ctx context.Context, settings map[string]string) error
}
