package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// TagsRepository 标签仓储接口
type TagsRepository interface {
	// GetAll 获取所有标签（含文章数量）
	GetAll(ctx context.Context) ([]dto.TagListItem, error)
}
