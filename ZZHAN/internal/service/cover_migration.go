package service

import (
	"ZZHAN/internal/model/entity"
	"ZZHAN/pkg/storage"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gabriel-vasile/mimetype"
	"gorm.io/gorm"
)

// MigrateArticleCovers replaces legacy data URLs after verifying the uploaded object.
func MigrateArticleCovers(ctx context.Context, db *gorm.DB, store storage.Storage) error {
	client := &http.Client{Timeout: 15 * time.Second}
	var lastID int64
	var failed []error
	for {
		var article entity.Article
		err := db.WithContext(ctx).Select("id", "cover_image").
			Where("id > ? AND cover_image LIKE ?", lastID, "data:image/%").
			Order("id").Take(&article).Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.Join(failed...)
		}
		if err != nil {
			return errors.Join(append(failed, err)...)
		}
		lastID = article.ID
		uploadCtx, cancel := context.WithTimeout(ctx, time.Minute)
		link, err := uploadLegacyCover(uploadCtx, client, store, article.ID, article.CoverImage)
		cancel()
		if err != nil {
			failed = append(failed, fmt.Errorf("文章 %d 封面迁移失败: %w", article.ID, err))
			continue
		}

		// Preserve the article if its cover was edited during the upload.
		oldHash := fmt.Sprintf("%x", sha256.Sum256([]byte(article.CoverImage)))
		result := db.WithContext(ctx).Model(&entity.Article{}).
			Where("id = ? AND SHA2(cover_image, 256) = ?", article.ID, oldHash).
			UpdateColumn("cover_image", link)
		if result.Error != nil {
			failed = append(failed, fmt.Errorf("文章 %d 封面链接保存失败: %w", article.ID, result.Error))
		}
	}
}

func uploadLegacyCover(ctx context.Context, client *http.Client, store storage.Storage, id int64, value string) (string, error) {
	header, encoded, ok := strings.Cut(value, ",")
	if !ok || !strings.HasSuffix(strings.ToLower(header), ";base64") {
		return "", fmt.Errorf("封面不是 Base64 图片")
	}
	data, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("Base64 解码失败: %w", err)
	}
	mime := mimetype.Detect(data).String()
	extension, ok := allowedImages[mime]
	if !ok {
		return "", fmt.Errorf("不支持的图片格式 %s", mime)
	}
	key := fmt.Sprintf("images/covers/%d-%x.%s", id, sha256.Sum256(data), extension)
	link, err := store.Upload(ctx, key, bytes.NewReader(data), int64(len(data)), mime)
	if err != nil {
		return "", err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Range", "bytes=0-0")
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if (resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent) ||
		!strings.HasPrefix(resp.Header.Get("Content-Type"), "image/") {
		return "", fmt.Errorf("图片链接不可读（HTTP %d）", resp.StatusCode)
	}
	return link, nil
}
