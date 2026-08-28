package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// TagsAdminService 后台标签业务接口
type TagsAdminService interface {
	// AdminList 获取标签列表
	AdminList(ctx context.Context, req *dto.AdminTagListRequest) (*dto.AdminTagListResponse, error)

	// AdminGetByID 获取标签详情
	AdminGetByID(ctx context.Context, id int64) (*dto.AdminTagItem, error)

	// AdminCreate 创建标签
	AdminCreate(ctx context.Context, req *dto.AdminTagCreateRequest) (*dto.AdminTagItem, error)

	// AdminUpdate 更新标签
	AdminUpdate(ctx context.Context, id int64, req *dto.AdminTagUpdateRequest) (*dto.AdminTagItem, error)

	// AdminUpdateStatus 修改标签状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) error

	// AdminDelete 删除标签
	AdminDelete(ctx context.Context, id int64) error
}
