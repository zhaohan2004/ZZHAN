package storage

import (
	"ZZHAN/pkg/config"
	"context"
	"fmt"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss"
	"github.com/aliyun/alibabacloud-oss-go-sdk-v2/oss/credentials"
	"io"
)

type ossStorage struct {
	client *oss.Client
	cfg    config.OSSStorageConfig
}

// objectURL 拼接对象访问 URL（固定默认域名）
func (s *ossStorage) objectURL(objectKey string) string {
	return fmt.Sprintf("https://%s.oss-%s.aliyuncs.com/%s", s.cfg.Bucket, s.cfg.Endpoint, objectKey)
}

// newOSSStorage 校验配置并创建 OSS 客户端。
func newOSSStorage(cfg config.OSSStorageConfig) (Storage, error) {
	if cfg.AccessKeyID == "" || cfg.AccessKeySecret == "" || cfg.Bucket == "" || cfg.Endpoint == "" {
		return nil, fmt.Errorf("oss 配置不完整：access_key_id/access_key_secret/bucket/endpoint 均必填")
	}
	ossCfg := oss.LoadDefaultConfig().
		WithCredentialsProvider(credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.AccessKeySecret, "")).
		WithRegion(cfg.Endpoint)
	return &ossStorage{client: oss.NewClient(ossCfg), cfg: cfg}, nil
}

// Upload 流式上传对象到 OSS。
// contentType 不传，OSS 会按 objectKey 扩展名自动识别 Content-Type。
func (s *ossStorage) Upload(ctx context.Context, objectKey string, r io.Reader, _ int64, _ string) (string, error) {
	request := &oss.PutObjectRequest{
		Bucket: oss.Ptr(s.cfg.Bucket),
		Key:    oss.Ptr(objectKey),
		Body:   r,
	}
	if _, err := s.client.PutObject(ctx, request); err != nil {
		return "", err
	}
	return s.objectURL(objectKey), nil
}
