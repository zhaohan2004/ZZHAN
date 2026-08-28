package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// TagsAdminRepository 后台标签数据访问接口
type TagsAdminRepository interface {
	// AdminList 分页查询标签列表
	AdminList(ctx context.Context, req *dto.AdminTagListRequest) ([]dto.AdminTagItem, int64, error)

	// AdminGetByID 通过 ID 获取标签
	AdminGetByID(ctx context.Context, id int64) (*dto.AdminTagItem, error)

	// AdminCreate 创建标签
	AdminCreate(ctx context.Context, name, status string) (int64, error)

	// AdminUpdate 更新标签
	AdminUpdate(ctx context.Context, id int64, name string) error

	// AdminUpdateStatus 修改标签状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) error

	// AdminDelete 删除标签
	AdminDelete(ctx context.Context, id int64) error
}
