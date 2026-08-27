package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// ArchivesRepository 归档仓储接口
type ArchivesRepository interface {
	// GetArchives 获取归档列表（按年月分组）
	GetArchives(ctx context.Context) ([]dto.ArchiveItem, error)
}
