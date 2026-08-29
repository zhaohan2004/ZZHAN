package service

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// OperationLogsAdminService 后台操作日志业务接口
type OperationLogsAdminService interface {
	// AdminList 获取操作日志列表
	AdminList(ctx context.Context, req *dto.OperationLogsAdminListRequest) (*dto.OperationLogsAdminListResponse, error)
}
