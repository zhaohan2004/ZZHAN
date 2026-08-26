package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// TagsService 标签业务接口
type TagsService interface {
	// GetAll 获取所有标签
	GetAll(ctx context.Context) ([]dto.TagListItem, error)
}
