package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type archivesService struct {
	archivesRepo repository.ArchivesRepository
}

// NewArchivesService 创建归档业务实例
func NewArchivesService(archivesRepo repository.ArchivesRepository) ArchivesService {
	return &archivesService{archivesRepo: archivesRepo}
}

// GetArchives 获取归档列表
func (s *archivesService) GetArchives(ctx context.Context) ([]dto.ArchiveItem, error) {
	return s.archivesRepo.GetArchives(ctx)
}
