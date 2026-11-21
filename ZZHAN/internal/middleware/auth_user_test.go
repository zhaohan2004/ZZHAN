package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"ZZHAN/internal/model/entity"
	"ZZHAN/internal/repository"
	"ZZHAN/pkg/config"
	appErrors "ZZHAN/pkg/errors"
	"ZZHAN/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

const authTestUserID = 42

type userAuthTestRepository struct {
	repository.AuthRepository
	user *entity.User
	err  error
}

func (r userAuthTestRepository) FindUserByID(_ context.Context, id int) (*entity.User, error) {
	if r.err != nil {
		return nil, r.err
	}
	if r.user == nil || int64(id) != r.user.ID {
		return nil, gorm.ErrRecordNotFound
	}
	return r.user, nil
}

type userAuthTestRedis struct {
	repository.RedisRepository
	available    bool
	activeTokens map[string]string
}

func (r userAuthTestRedis) Available(context.Context) bool {
	return r.available
}

func (r userAuthTestRedis) IsBlacklisted(context.Context, string) (bool, error) {
	return false, nil
}

func (r userAuthTestRedis) GetActiveToken(_ context.Context, scope string, id int) (string, error) {
	token, found := r.activeTokens[scope]
	if !found || id != authTestUserID {
		return "", redis.Nil
	}
	return token, nil
}

func TestUserAuthAccountStatus(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-signing-key-only-for-user-auth-unit-tests")
	t.Setenv("ZZHAN_JWT_ACCESS_EXPIRE_HOURS", "")
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	if err := os.WriteFile(configPath, []byte("jwt:\n  access_expire_hours: 1\n"), 0600); err != nil {
		t.Fatalf("write test config: %v", err)
	}
	if _, err := config.Load(configPath); err != nil {
		t.Fatalf("load test config: %v", err)
	}
	legacyToken, err := jwt.GenerateToken(authTestUserID, "reader", "access_token")
	if err != nil {
		t.Fatalf("sign legacy test token: %v", err)
	}
	currentToken, err := jwt.GenerateUserToken(authTestUserID, "reader", "access_token", 2)
	if err != nil {
		t.Fatalf("sign current test token: %v", err)
	}

	activeTokens := map[string]string{"user": legacyToken, "user:2": currentToken}
	for _, tt := range []struct {
		name        string
		token       string
		status      int8
		version     uint64
		redisOnline bool
		missing     bool
		databaseErr error
		wantStatus  int
	}{
		{name: "legacy session remains valid", token: legacyToken, status: 1, redisOnline: true, wantStatus: http.StatusNoContent},
		{name: "disabled with cached session", token: legacyToken, status: 0, version: 1, redisOnline: true, wantStatus: http.StatusForbidden},
		{name: "disabled while Redis is down", token: legacyToken, status: 0, version: 1, wantStatus: http.StatusForbidden},
		{name: "old session stays revoked after reenabling", token: legacyToken, status: 1, version: 2, redisOnline: true, wantStatus: http.StatusUnauthorized},
		{name: "old session stays revoked without Redis", token: legacyToken, status: 1, version: 2, wantStatus: http.StatusUnauthorized},
		{name: "current session succeeds after reenabling", token: currentToken, status: 1, version: 2, redisOnline: true, wantStatus: http.StatusNoContent},
		{name: "deleted user", token: legacyToken, missing: true, redisOnline: true, wantStatus: http.StatusUnauthorized},
		{name: "database outage", token: legacyToken, databaseErr: errors.New("database unavailable"), redisOnline: true, wantStatus: http.StatusServiceUnavailable},
		{name: "no token", status: 1, redisOnline: true, wantStatus: http.StatusUnauthorized},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var user *entity.User
			if !tt.missing {
				user = &entity.User{
					BaseEntity:     entity.BaseEntity{ID: authTestUserID},
					Status:         tt.status,
					SessionVersion: tt.version,
				}
			}
			authRepo := userAuthTestRepository{user: user, err: tt.databaseErr}
			redisRepo := userAuthTestRedis{available: tt.redisOnline, activeTokens: activeTokens}
			authenticated := tt.wantStatus == http.StatusNoContent

			t.Run("protected route", func(t *testing.T) {
				checkUserAuthRequest(t, AuthUser(redisRepo, authRepo), tt.token, tt.wantStatus, authenticated)
			})
			t.Run("public route", func(t *testing.T) {
				checkUserAuthRequest(t, OptionalUserAuth(redisRepo, authRepo), tt.token, http.StatusNoContent, authenticated)
			})
		})
	}
}

func checkUserAuthRequest(t *testing.T, auth gin.HandlerFunc, token string, wantStatus int, wantAuthenticated bool) {
	t.Helper()
	var called bool
	var userID int
	var username string
	router := gin.New()
	router.GET("/", auth, func(c *gin.Context) {
		called = true
		userID = GetUserID(c)
		username = GetUsername(c)
		c.Status(http.StatusNoContent)
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	res := httptest.NewRecorder()
	router.ServeHTTP(res, req)

	if res.Code != wantStatus || called != (wantStatus == http.StatusNoContent) {
		t.Fatalf("got status %d and downstream handler called=%t, want status %d", res.Code, called, wantStatus)
	}
	if wantAuthenticated && (userID != authTestUserID || username != "reader") {
		t.Errorf("active user identity was not populated")
	}
	if !wantAuthenticated && (userID != 0 || username != "") {
		t.Errorf("anonymous request received a user identity")
	}
	if wantStatus == http.StatusForbidden {
		var response struct {
			Code int `json:"code"`
		}
		if err := json.Unmarshal(res.Body.Bytes(), &response); err != nil {
			t.Fatalf("decode disabled-user response: %v", err)
		}
		if response.Code != appErrors.CodeUserDisabled {
			t.Errorf("got business code %d, want %d", response.Code, appErrors.CodeUserDisabled)
		}
	}
}
