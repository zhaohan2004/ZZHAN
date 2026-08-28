package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// ArticlesAdminRepository 后台文章仓储接口
type ArticlesAdminRepository interface {
	// AdminList 获取文章列表（全部状态，支持筛选）
	AdminList(ctx context.Context, req *dto.AdminArticleListRequest) ([]dto.AdminArticleItem, int64, error)

	// AdminGetByID 通过 ID 获取文章详情
	AdminGetByID(ctx context.Context, id int64) (*dto.AdminArticleDetail, error)

	// AdminCreate 创建文章
	AdminCreate(ctx context.Context, req *dto.AdminArticleCreateRequest, authorID int64) (int64, error)

	// AdminUpdate 更新文章
	AdminUpdate(ctx context.Context, id int64, req *dto.AdminArticleUpdateRequest) error

	// AdminDelete 删除文章（软删除）
	AdminDelete(ctx context.Context, id int64) error

	// AdminUpdateStatus 修改文章状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) error
}
