package entity

// Like 文章点赞
type Like struct {
	CreatedOnlyEntity
	ArticleID uint64 `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_likes,priority:1;comment:文章ID" json:"article_id"`
	UserID    uint64 `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_likes,priority:2;comment:用户ID" json:"user_id"`
}

// TableName 指定表名
func (Like) TableName() string {
	return "likes"
}
