package entity

import (
	"time"

	"gorm.io/gorm"
)

// Article 文章
type Article struct {
	BaseEntity
	Title        string         `gorm:"type:varchar(200);not null;comment:文章标题" json:"title"`
	Slug         string         `gorm:"type:varchar(220);not null;uniqueIndex;comment:SEO 友好 URL" json:"slug"`
	Summary      string         `gorm:"type:varchar(500);comment:摘要，用于列表与 SEO" json:"summary"`
	CoverImage   string         `gorm:"type:varchar(255);comment:封面图 URL（空则前端渐变生成）" json:"coverImage"`
	CategoryID   uint64         `gorm:"type:bigint unsigned;not null;index;comment:所属分类ID" json:"categoryId"`
	AuthorID     uint64         `gorm:"type:bigint unsigned;not null;index;comment:作者 = 后台管理员ID" json:"authorId"`
	Content      string         `gorm:"type:longtext;not null;comment:Markdown 正文" json:"content"`
	Status       string         `gorm:"type:varchar(20);not null;default:draft;index:idx_articles_status_published,priority:1;comment:状态：draft / published / down" json:"status"`
	Featured     int8           `gorm:"type:tinyint;not null;default:0;comment:是否编辑推荐：1 是 / 0 否" json:"featured"`
	Hot          int8           `gorm:"type:tinyint;not null;default:0;comment:是否热门：1 是 / 0 否" json:"hot"`
	Views        uint           `gorm:"type:int unsigned;not null;default:0;index;comment:浏览量(冗余计数)" json:"views"`
	Likes        uint           `gorm:"type:int unsigned;not null;default:0;comment:点赞数(冗余计数)" json:"likes"`
	CommentCount uint           `gorm:"type:int unsigned;not null;default:0;comment:评论数(冗余计数)" json:"commentCount"`
	PublishedAt  *time.Time     `gorm:"type:datetime;index:idx_articles_status_published,priority:2;comment:发布时间" json:"publishedAt"`
	DeletedAt    gorm.DeletedAt `gorm:"type:datetime;index;comment:软删除时间(非空=已删除)" json:"-"`
}

// TableName 指定表名
func (Article) TableName() string {
	return "articles"
}
