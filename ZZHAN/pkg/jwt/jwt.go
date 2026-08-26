package jwt

import (
	"errors"
	"time"

	"ZZHAN/pkg/config"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrTokenInvalid = errors.New("token 无效")
	ErrTokenExpired = errors.New("token 已过期")
)

// GenerateToken 生成指定类型的 JWT Token
func GenerateToken(userID int, username string, tokenType string) (string, error) {
	cfg := config.Get().JWT

	var expireDuration time.Duration
	if tokenType == "refresh_token" {
		expireDuration = cfg.RefreshExpireHours * time.Hour
	} else {
		expireDuration = cfg.AccessExpireHours * time.Hour
	}

	claims := CustomClaims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expireDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "ZZHAN",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// GenerateTokenPair 生成access + refresh token 对
func GenerateTokenPair(userID int, username string) (accessToken, refreshToken string, err error) {
	accessToken, err = GenerateToken(userID, username, "access_token")
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateToken(userID, username, "refresh_token")
	if err != nil {
		return "", "", err
	}
	return accessToken, refreshToken, nil
}

// ParseToken 解析 JWT Token
func ParseToken(tokenString string, expectedType string) (*CustomClaims, error) {
	cfg := config.Get().JWT

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(cfg.Secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if expectedType != "" && claims.TokenType != expectedType {
			return nil, ErrTokenInvalid
		}
		return claims, nil
	}

	return nil, ErrTokenInvalid
}
