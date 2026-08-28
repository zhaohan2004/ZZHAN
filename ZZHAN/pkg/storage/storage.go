package storage

import (
	"ZZHAN/pkg/config"
	"context"
	"fmt"
	"io"
)

type Storage interface {
	// Upload 将文件内容流式上传到指定对象键（objectKey），返回可访问的永久 URL。
	// objectKey 形如 "images/2026/08/08/abc123.jpg"，由调用方生成。
	Upload(ctx context.Context, objectKey string, r io.Reader, size int64, contentType string) (string, error)
}

// New 根据配置创建存储实现。
func New(cfg config.StorageConfig) (Storage, error) {
	switch cfg.Driver {
	case "oss":
		return newOSSStorage(cfg.OSS)
	case "local":
		return newLocalStorage(cfg.Local)
	default:
		return nil, fmt.Errorf("不支持的存储驱动")
	}
}
