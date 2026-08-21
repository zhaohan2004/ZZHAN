package entity

// ArticleTag 文章-标签关联表（多对多）
type ArticleTag struct {
	ArticleID uint64 `gorm:"type:bigint unsigned;primaryKey;comment:文章ID" json:"articleId"`
	TagID     uint64 `gorm:"type:bigint unsigned;primaryKey;index;comment:标签ID" json:"tagId"`
}

// TableName 指定表名
func (ArticleTag) TableName() string {
	return "article_tags"
}
