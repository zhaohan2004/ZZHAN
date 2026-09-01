package service

import (
	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"
	"ZZHAN/internal/repository"
	"ZZHAN/pkg/config"
	"ZZHAN/pkg/jwt"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

// githubHTTPClient 根据配置创建用于请求 GitHub API 的 HTTP 客户端
func githubHTTPClient() *http.Client {
	proxy := config.Get().GitHub.Proxy
	if proxy == "" {
		return &http.Client{Timeout: 15 * time.Second}
	}
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return &http.Client{Timeout: 15 * time.Second}
	}
	return &http.Client{
		Timeout: 15 * time.Second,
		Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
		},
	}
}

type authService struct {
	authRepo  repository.AuthRepository
	redisRepo repository.RedisRepository
}

func NewAuthService(authRepo repository.AuthRepository, redisRepo repository.RedisRepository) AuthService {
	return &authService{authRepo: authRepo, redisRepo: redisRepo}
}

// GitHubUser GitHub 用户信息
type GitHubUser struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
}

// GitHubToken GitHub access_token 响应
type GitHubToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

// GitHubLogin GitHub OAuth 登录
func (s *authService) GitHubLogin(ctx context.Context, code, redirectURI string) (*dto.LoginResponse, error) {
	cfg := config.Get().GitHub

	//1.先用 code 换取 GitHub 的 access_token
	ghToken, err := s.exchangeGitHubToken(cfg.ClientID, cfg.ClientSecret, code, redirectURI)
	if err != nil {
		return nil, fmt.Errorf("获取 GitHub access_token 失败：%w", err)
	}

	//2.再用GitHub access_token 获取用户信息
	ghUser, err := s.getGitHubUser(ghToken.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("获取 GitHub 用户信息失败：%w", err)
	}

	//3.然后查找/创建用户
	user, err := s.findOrCreateUser(ctx, ghUser)
	if err != nil {
		return nil, fmt.Errorf("查找/创建用户失败：%w", err)
	}

	//4.异步更新登录时间

	go func() {
		_ = s.authRepo.UpdateUserLoginTime(context.Background(), int(user.ID))
	}()

	//5.生成access_token 和 refresh_token 的 JWT token 对
	accessToken, refreshToken, err := jwt.GenerateTokenPair(int(user.ID), user.Nickname)
	if err != nil {
		return nil, fmt.Errorf("生成 token 失败：%w", err)
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.AuthUser{
			ID:       int(user.ID),
			Provider: user.Provider,
			Nickname: user.Nickname,
			Avatar:   user.Avatar,
		},
	}, err
}

// exchangeGitHubToken 用 code 换取 GitHub access_token
func (s *authService) exchangeGitHubToken(clientID, clientSecret, code, redirectURI string) (*GitHubToken, error) {
	url := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	if redirectURI != "" {
		url += "&redirect_uri=" + redirectURI
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := githubHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var token GitHubToken
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("GitHub 返回的 access_token 为空")
	}
	return &token, nil
}

// getGitHubUser 获取 GitHub 用户信息
func (s *authService) getGitHubUser(accessToken string) (*GitHubUser, error) {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := githubHTTPClient().Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var user GitHubUser
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// findOrCreateUser 查找或创建用户
func (s *authService) findOrCreateUser(ctx context.Context, ghUser *GitHubUser) (*entity.User, error) {
	user, err := s.authRepo.FindUserByOpenID(ctx, "github", fmt.Sprintf("%d", ghUser.ID))
	if err == nil {
		//用户存在，获取最新的头像和昵称
		user.Avatar = ghUser.AvatarURL
		if ghUser.Name != "" {
			user.Nickname = ghUser.Name
		}
		_ = s.authRepo.UpdateUser(ctx, user)
		return user, nil
	}

	//用户不存在在时，创建新用户
	nickname := ghUser.Name
	if nickname == "" {
		nickname = ghUser.Login
	}
	newUser := entity.User{
		Provider: "github",
		Openid:   fmt.Sprintf("%d", ghUser.ID),
		Nickname: nickname,
		Avatar:   ghUser.AvatarURL,
		Status:   1,
	}
	if err := s.authRepo.CreateUser(ctx, &newUser); err != nil {
		return nil, err
	}
	return &newUser, nil
}

// RefreshToken 刷新 access_token
func (s *authService) RefreshToken(ctx context.Context, refreshToken string) (*dto.RefreshResponse, error) {
	// 验证 refresh_token
	claims, err := jwt.ParseToken(refreshToken, "refresh_token")
	if err != nil {
		return nil, fmt.Errorf("refresh_token 无效：%w", err)
	}

	// 检查是否在黑名单中
	blacklisted, err := s.redisRepo.IsBlacklisted(ctx, refreshToken)
	if err != nil {
		return nil, fmt.Errorf("检查 token 黑名单失败：%w", err)
	}
	if blacklisted {
		return nil, fmt.Errorf("refresh_token 已失效")
	}

	newAccessToken, err := jwt.GenerateToken(claims.UserID, claims.Username, "access_token")
	if err != nil {
		return nil, fmt.Errorf("生成新的 access_token 失败：%w", err)
	}
	return &dto.RefreshResponse{
		AccessToken: newAccessToken,
	}, nil
}

// Logout 退出登录
func (s *authService) Logout(ctx context.Context, accessToken string) error {
	claims, err := jwt.ParseToken(accessToken, "access_token")
	if err != nil {
		//token 已失效，可以直接返回成功
		return nil
	}

	// 计算剩余过期时间
	expireAt := claims.ExpiresAt.Time
	remaining := time.Until(expireAt)
	if remaining <= 0 {
		//token 已过期，直接返回成功
		return nil
	}

	// 加入黑名单
	return s.redisRepo.AddToBlacklist(ctx, accessToken, remaining)
}

// GetCurrentUser 获取当前登录用户信息
func (s *authService) GetCurrentUser(ctx context.Context, userID int) (*dto.AuthUser, error) {
	user, err := s.authRepo.FindUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("用户不存在：%w", err)
	}

	return &dto.AuthUser{
		ID:       int(user.ID),
		Provider: user.Provider,
		Nickname: user.Nickname,
		Avatar:   user.Avatar,
	}, nil
}
