package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// ArticlesAdminService 后台文章业务接口
type ArticlesAdminService interface {
	// AdminList 获取文章列表
	AdminList(ctx context.Context, req *dto.AdminArticleListRequest) (*dto.AdminArticleListResponse, error)

	// AdminGetByID 获取文章详情
	AdminGetByID(ctx context.Context, id int64) (*dto.AdminArticleDetail, error)

	// AdminCreate 创建文章
	AdminCreate(ctx context.Context, req *dto.AdminArticleCreateRequest, authorID int64) (*dto.AdminArticleDetail, error)

	// AdminUpdate 更新文章
	AdminUpdate(ctx context.Context, id int64, req *dto.AdminArticleUpdateRequest) (*dto.AdminArticleDetail, error)

	// AdminDelete 删除文章
	AdminDelete(ctx context.Context, id int64) error

	// AdminUpdateStatus 修改文章状态
	AdminUpdateStatus(ctx context.Context, id int64, status string) (*dto.AdminArticleDetail, error)
}
