package storage

import (
	"ZZHAN/pkg/config"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// localStorage 本地磁盘存储实现（driver=local）。
// 文件写 {baseDir}/{objectKey}，访问 URL = {baseURL}/static/{objectKey}（静态路由在 app 层挂载）。
type localStorage struct {
	baseDir string
	baseURL string
}

// newLocalStorage 校验配置并创建本地存储。
func newLocalStorage(cfg config.LocalStorageConfig) (Storage, error) {
	if cfg.BaseDir == "" || cfg.BaseURL == "" {
		return nil, fmt.Errorf("local 配置不完整：base_dir/base_url 均必填")
	}
	if err := os.MkdirAll(cfg.BaseDir, 0o755); err != nil {
		return nil, fmt.Errorf("创建本地存储目录失败: %w", err)
	}
	return &localStorage{
		baseDir: cfg.BaseDir,
		baseURL: strings.TrimRight(cfg.BaseURL, "/"),
	}, nil
}

// Upload 流式写文件到本地磁盘（自动建目录），返回可访问 URL。
// contentType 忽略：静态文件按扩展名由浏览器识别。
func (s *localStorage) Upload(_ context.Context, objectKey string, r io.Reader, _ int64, _ string) (string, error) {
	fullPath := filepath.Join(s.baseDir, filepath.FromSlash(objectKey))
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", err
	}
	dst, err := os.Create(fullPath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	if _, err := io.Copy(dst, r); err != nil {
		return "", err
	}
	return s.objectURL(objectKey), nil
}

// objectURL 拼接本地文件访问 URL
func (s *localStorage) objectURL(objectKey string) string {
	return s.baseURL + "/static/" + objectKey
}
