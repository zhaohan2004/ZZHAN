package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CommentsAdminRepository 后台评论仓储接口
type CommentsAdminRepository interface {
	// AdminList 获取评论列表（支持筛选）
	AdminList(ctx context.Context, req *dto.CommentsAdminListRequest) ([]dto.CommentsAdminItem, int64, error)

	// AdminGetByID 通过 ID 获取评论详情
	AdminGetByID(ctx context.Context, id int64) (*dto.CommentsAdminItem, error)

	// AdminUpdateStatus 修改评论状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) error

	// AdminDelete 删除评论（软删除）
	AdminDelete(ctx context.Context, id int64) error
}
