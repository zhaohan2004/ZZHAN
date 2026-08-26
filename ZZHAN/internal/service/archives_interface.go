package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// ArchivesService 归档业务接口
type ArchivesService interface {
	// GetArchives 获取归档列表
	GetArchives(ctx context.Context) ([]dto.ArchiveItem, error)
}
