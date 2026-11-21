package service

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"ZZHAN/internal/model/entity"
	"ZZHAN/internal/repository"
	"ZZHAN/pkg/config"
	appErrors "ZZHAN/pkg/errors"
	"ZZHAN/pkg/jwt"
)

type loginTestRepository struct {
	repository.AuthRepository
	user entity.User
}

func (r *loginTestRepository) FindUserByOpenID(context.Context, string, string) (*entity.User, error) {
	user := r.user
	return &user, nil
}

func (r *loginTestRepository) FindUserByID(context.Context, int) (*entity.User, error) {
	user := r.user
	return &user, nil
}

func (r *loginTestRepository) UpdateUser(context.Context, *entity.User) error {
	return nil
}

func (r *loginTestRepository) UpdateUserLoginTime(context.Context, int) error {
	return nil
}

type loginTestTransport func(*http.Request) (*http.Response, error)

func (f loginTestTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestLoginAndRefreshRespectAccountStatus(t *testing.T) {
	t.Setenv("JWT_SECRET", "test-key-for-login-and-refresh")
	t.Setenv("ZZHAN_GITHUB_PROXY", "")
	path := filepath.Join(t.TempDir(), "config.yaml")
	settings := []byte("jwt:\n  access_expire_hours: 1\n  refresh_expire_hours: 2\ngithub:\n  proxy: ''\n")
	if err := os.WriteFile(path, settings, 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := config.Load(path); err != nil {
		t.Fatal(err)
	}

	originalTransport := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = originalTransport })
	http.DefaultTransport = loginTestTransport(func(req *http.Request) (*http.Response, error) {
		var body string
		switch req.URL.Host {
		case "github.com":
			body = `{"access_token":"test-github-token"}`
		case "api.github.com":
			body = `{"id":123,"login":"reader"}`
		default:
			t.Fatalf("unexpected OAuth URL: %s", req.URL)
		}
		return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}, nil
	})

	repo := &loginTestRepository{user: entity.User{
		BaseEntity: entity.BaseEntity{ID: 2},
		Status:     1,
	}}
	auth := NewAuthService(repo, nil)
	ctx := t.Context()
	initial, err := auth.GitHubLogin(ctx, "test-code", "")
	if err != nil {
		t.Fatalf("initial login: %v", err)
	}
	if _, err := auth.RefreshToken(ctx, initial.RefreshToken); err != nil {
		t.Fatalf("initial refresh: %v", err)
	}

	repo.user.Status = 0
	repo.user.SessionVersion = 1
	if response, err := auth.GitHubLogin(ctx, "test-code", ""); response != nil || !errors.Is(err, appErrors.ErrUserDisabled) {
		t.Fatalf("login while disabled returned a response or unexpected error: %v", err)
	}
	if response, err := auth.RefreshToken(ctx, initial.RefreshToken); response != nil || !errors.Is(err, appErrors.ErrUserDisabled) {
		t.Fatalf("refresh while disabled returned a response or unexpected error: %v", err)
	}

	repo.user.Status = 1
	repo.user.SessionVersion = 2
	if response, err := auth.RefreshToken(ctx, initial.RefreshToken); response != nil || !errors.Is(err, appErrors.ErrInvalidToken) {
		t.Fatalf("refresh from an earlier session returned a response or unexpected error: %v", err)
	}
	current, err := auth.GitHubLogin(ctx, "test-code", "")
	if err != nil {
		t.Fatalf("login after reenabling: %v", err)
	}
	claims, err := jwt.ParseToken(current.AccessToken, "access_token")
	if err != nil || claims.SessionVersion != 2 {
		t.Fatalf("login did not issue a token for the current session: %v", err)
	}
	renewed, err := auth.RefreshToken(ctx, current.RefreshToken)
	if err != nil {
		t.Fatalf("refresh after reenabling: %v", err)
	}
	claims, err = jwt.ParseToken(renewed.AccessToken, "access_token")
	if err != nil || claims.SessionVersion != 2 {
		t.Fatalf("refresh did not preserve the current session version: %v", err)
	}
}
