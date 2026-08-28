package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CategoriesAdminRepository 后台分类数据访问接口
type CategoriesAdminRepository interface {
	// AdminList 分页查询分类列表
	AdminList(ctx context.Context, req *dto.AdminCategoryListRequest) ([]dto.AdminCategoryItem, int64, error)

	// AdminGetByID 通过 ID 获取分类
	AdminGetByID(ctx context.Context, id int64) (*dto.AdminCategoryItem, error)

	// AdminCreate 创建分类
	AdminCreate(ctx context.Context, name, slug, desc, color, status string) (int64, error)

	// AdminUpdate 更新分类
	AdminUpdate(ctx context.Context, id int64, updates map[string]interface{}) error

	// AdminUpdateStatus 修改分类状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) error

	// AdminDelete 删除分类
	AdminDelete(ctx context.Context, id int64) error
}
