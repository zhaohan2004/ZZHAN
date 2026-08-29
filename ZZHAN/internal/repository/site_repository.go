package repository

import (
	"context"
	"encoding/json"
	"strconv"

	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type siteRepository struct {
	db *gorm.DB
}

// NewSiteRepository 创建站点设置仓储
func NewSiteRepository(db *gorm.DB) SiteRepository {
	return &siteRepository{db: db}
}

// GetSiteResponse 从 site_settings 表读取所有 KV，组装为 SiteResponse
func (r *siteRepository) GetSiteResponse(ctx context.Context) (*dto.SiteResponse, error) {
	settings, err := r.GetSettings(ctx)
	if err != nil {
		return nil, err
	}

	resp := &dto.SiteResponse{
		Name:         settings["blog_name"],
		Tagline:      pickFirst(settings["tagline"], settings["blog_desc"]),
		Bio:          settings["author_intro"],
		Github:       settings["github"],
		Email:        settings["email"],
		Author:       settings["author_name"],
		Role:         settings["author_role"],
		Motto:        settings["motto"],
		Location:     settings["location"],
		Avatar:       settings["avatar"],
		HeroTerminal: settings["hero_terminal"],
	}

	// since: string -> int
	if v, err := strconv.Atoi(settings["since"]); err == nil {
		resp.Since = v
	}

	// socials: JSON string -> []SocialsList
	if raw := settings["socials"]; raw != "" {
		var socials []dto.SocialsList
		if err := json.Unmarshal([]byte(raw), &socials); err == nil {
			resp.Socials = socials
		}
	}

	return resp, nil
}

// GetSettings 查询 site_settings 全表，返回 key -> value map
func (r *siteRepository) GetSettings(ctx context.Context) (map[string]string, error) {
	var rows []entity.SiteSetting
	if err := r.db.WithContext(ctx).Find(&rows).Error; err != nil {
		return nil, err
	}
	result := make(map[string]string, len(rows))
	for _, row := range rows {
		result[row.Key] = row.Value
	}
	return result, nil
}

// UpdateSettings 批量 upsert 设置项
func (r *siteRepository) UpdateSettings(ctx context.Context, settings map[string]string) error {
	if len(settings) == 0 {
		return nil
	}
	rows := make([]entity.SiteSetting, 0, len(settings))
	for k, v := range settings {
		rows = append(rows, entity.SiteSetting{
			Key:   k,
			Value: v,
		})
	}
	// ON DUPLICATE KEY UPDATE value
	return r.db.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "key"}},
			DoUpdates: clause.AssignmentColumns([]string{"value"}),
		}).
		Create(&rows).Error
}

// pickFirst 返回第一个非空字符串
func pickFirst(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}
