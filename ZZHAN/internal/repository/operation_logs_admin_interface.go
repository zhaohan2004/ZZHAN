package repository

import (
	"ZZHAN/internal/model/dto"
	"context"
)

// OperationLogsAdminRepository 后台操作日志仓储接口
type OperationLogsAdminRepository interface {
	// AdminList 获取操作日志列表（支持筛选）
	AdminList(ctx context.Context, req *dto.OperationLogsAdminListRequest) ([]dto.OperationLogsAdminItem, int64, error)
}
