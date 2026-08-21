package service

import (
	"context"

	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
)

type siteService struct {
	siteRepo repository.SiteRepository
}

// NewSiteService 创建站点设置业务实例
func NewSiteService(siteRepo repository.SiteRepository) SiteService {
	return &siteService{siteRepo: siteRepo}
}

// GetSite 前台获取站点信息 — 直接委托 repository 组装 SiteResponse
func (s *siteService) GetSite(ctx context.Context) (*dto.SiteResponse, error) {
	return s.siteRepo.GetSiteResponse(ctx)
}

// GetSettings 后台获取全部设置 KV
func (s *siteService) GetSettings(ctx context.Context) (map[string]string, error) {
	return s.siteRepo.GetSettings(ctx)
}

// UpdateSettings 后台保存设置
func (s *siteService) UpdateSettings(ctx context.Context, settings map[string]string) error {
	return s.siteRepo.UpdateSettings(ctx, settings)
}
