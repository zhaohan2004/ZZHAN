package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CommentsAdminService 后台评论业务接口
type CommentsAdminService interface {
	// AdminList 获取评论列表
	AdminList(ctx context.Context, req *dto.CommentsAdminListRequest) (*dto.CommentsAdminListResponse, error)

	// AdminGetByID 获取评论详情
	AdminGetByID(ctx context.Context, id int64) (*dto.CommentsAdminItem, error)

	// AdminUpdateStatus 修改评论状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) (*dto.CommentsAdminItem, error)

	// AdminDelete 删除评论
	AdminDelete(ctx context.Context, id int64) error
}
