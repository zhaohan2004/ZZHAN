package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/repository"
	"context"
)

type operationLogsAdminService struct {
	operationLogsAdminRepo repository.OperationLogsAdminRepository
}

// NewOperationLogsAdminService 创建后台操作日志业务实例
func NewOperationLogsAdminService(operationLogsAdminRepo repository.OperationLogsAdminRepository) OperationLogsAdminService {
	return &operationLogsAdminService{operationLogsAdminRepo: operationLogsAdminRepo}
}

// AdminList 获取操作日志列表
func (s *operationLogsAdminService) AdminList(ctx context.Context, req *dto.OperationLogsAdminListRequest) (*dto.OperationLogsAdminListResponse, error) {
	// 分页默认值
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	logs, total, err := s.operationLogsAdminRepo.AdminList(ctx, req)
	if err != nil {
		return nil, err
	}

	return &dto.OperationLogsAdminListResponse{
		List:  logs,
		Total: total,
	}, nil
}
