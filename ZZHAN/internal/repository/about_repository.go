package repository

import (
	"context"
	"encoding/json"

	"ZZHAN/internal/model/dto"
	"ZZHAN/internal/model/entity"

	"gorm.io/gorm"
)

type aboutRepository struct {
	db *gorm.DB
}

// NewAboutRepository 创建关于我仓储实例
func NewAboutRepository(db *gorm.DB) AboutRepository {
	return &aboutRepository{db: db}
}

// GetAbout 获取关于我数据
func (r *aboutRepository) GetAbout(ctx context.Context) (*dto.AboutResponse, error) {
	resp := &dto.AboutResponse{
		Skills: []dto.Skill{},
	}

	// 查询 skills 分区的条目，按排序值升序
	var items []entity.AboutItem
	if err := r.db.WithContext(ctx).
		Where("section = ?", "skills").
		Order("sort_order ASC, id ASC").
		Find(&items).Error; err != nil {
		return nil, err
	}

	// 解析每条记录：Title -> name, Content(JSON) -> level
	for _, item := range items {
		var content struct {
			Level int `json:"level"`
		}
		if err := json.Unmarshal([]byte(item.Content), &content); err != nil {
			continue // 跳过格式异常的记录
		}

		name := ""
		if item.Title != nil {
			name = *item.Title
		}

		resp.Skills = append(resp.Skills, dto.Skill{
			Name:  name,
			Level: content.Level,
		})
	}

	return resp, nil
}
