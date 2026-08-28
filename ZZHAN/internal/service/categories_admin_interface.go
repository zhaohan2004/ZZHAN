package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// CategoriesAdminService 后台分类业务接口
type CategoriesAdminService interface {
	// AdminList 获取分类列表
	AdminList(ctx context.Context, req *dto.AdminCategoryListRequest) (*dto.AdminCategoryListResponse, error)

	// AdminGetByID 获取分类详情
	AdminGetByID(ctx context.Context, id int64) (*dto.AdminCategoryItem, error)

	// AdminCreate 创建分类
	AdminCreate(ctx context.Context, req *dto.AdminCategoryCreateRequest) (*dto.AdminCategoryItem, error)

	// AdminUpdate 更新分类
	AdminUpdate(ctx context.Context, id int64, req *dto.AdminCategoryUpdateRequest) (*dto.AdminCategoryItem, error)

	// AdminUpdateStatus 修改分类状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) error

	// AdminDelete 删除分类
	AdminDelete(ctx context.Context, id int64) error
}
