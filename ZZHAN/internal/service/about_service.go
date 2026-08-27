package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type aboutService struct {
	aboutRepo repository.AboutRepository
}

// NewAboutService 创建关于我业务实例
func NewAboutService(aboutRepo repository.AboutRepository) AboutService {
	return &aboutService{aboutRepo: aboutRepo}
}

// GetAbout 获取关于我数据
func (s *aboutService) GetAbout(ctx context.Context) (*dto.AboutResponse, error) {
	return s.aboutRepo.GetAbout(ctx)
}
