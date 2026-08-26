package repository

import (
	"context"
	"time"

	"ZZHAN/internal/model/entity"

	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

// NewAuthRepository 创建认证仓储
func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{db: db}
}

// FindUserByOpenID 根据 provider 和 openid 查找用户
func (r *authRepository) FindUserByOpenID(ctx context.Context, provider, openid string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("provider = ? AND openid = ?", provider, openid).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindUserByID 根据用户ID查找用户
func (r *authRepository) FindUserByID(ctx context.Context, userID int) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (r *authRepository) CreateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// UpdateUser 更新用户信息
func (r *authRepository) UpdateUser(ctx context.Context, user *entity.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// UpdateUserLoginTime 更新用户最后登录时间
func (r *authRepository) UpdateUserLoginTime(ctx context.Context, userID int) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&entity.User{}).
		Where("id = ?", userID).
		Update("last_login_at", now).Error
}
