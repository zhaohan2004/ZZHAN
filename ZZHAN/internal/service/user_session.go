package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"ZZHAN/internal/repository"
	appErrors "ZZHAN/pkg/errors"

	"gorm.io/gorm"
)

// ValidateUserSession checks the account and session against their current database state.
func ValidateUserSession(ctx context.Context, authRepo repository.AuthRepository, userID int, sessionVersion uint64) error {
	user, err := authRepo.FindUserByID(ctx, userID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return appErrors.ErrInvalidToken
	}
	if err != nil {
		return fmt.Errorf("%w: 读取用户状态失败: %w", appErrors.ErrInternalError, err)
	}
	if user.Status != 1 {
		return appErrors.ErrUserDisabled
	}
	if user.SessionVersion != sessionVersion {
		return appErrors.ErrInvalidToken
	}
	return nil
}

// UserSessionScope keeps requests from earlier versions out of the current Redis session.
func UserSessionScope(sessionVersion uint64) string {
	if sessionVersion == 0 {
		return "user"
	}
	return "user:" + strconv.FormatUint(sessionVersion, 10)
}
