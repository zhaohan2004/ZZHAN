package service

import (
	"context"

	"ZZHAN/internal/model/dto"
)

// SiteService 站点设置业务接口
type SiteService interface {
	// GetSite 前台获取站点信息
	GetSite(ctx context.Context) (*dto.SiteResponse, error)

	// GetSettings 后台获取全部设置
	GetSettings(ctx context.Context) (map[string]string, error)

	// UpdateSettings 后台保存设置
	UpdateSettings(ctx context.Context, settings map[string]string) error
}
