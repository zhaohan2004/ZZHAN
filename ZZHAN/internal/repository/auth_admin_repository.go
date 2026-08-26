package repository

import (
	"context"

	"ZZHAN/internal/model/entity"

	"gorm.io/gorm"
)

type adminAuthRepository struct {
	db *gorm.DB
}

// NewAdminAuthRepository 创建后台认证仓储
func NewAdminAuthRepository(db *gorm.DB) AdminAuthRepository {
	return &adminAuthRepository{db: db}
}

// FindAdminByUsername 根据用户名查找管理员
func (r *adminAuthRepository) FindAdminByUsername(ctx context.Context, username string) (*entity.Admin, error) {
	var admin entity.Admin
	err := r.db.WithContext(ctx).
		Where("username = ?", username).
		First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// FindAdminByID 根据 ID 查找管理员
func (r *adminAuthRepository) FindAdminByID(ctx context.Context, id int) (*entity.Admin, error) {
	var admin entity.Admin
	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&admin).Error
	if err != nil {
		return nil, err
	}
	return &admin, nil
}

// UpdateAdmin 更新管理员信息
func (r *adminAuthRepository) UpdateAdmin(ctx context.Context, admin *entity.Admin) error {
	return r.db.WithContext(ctx).Save(admin).Error
}
