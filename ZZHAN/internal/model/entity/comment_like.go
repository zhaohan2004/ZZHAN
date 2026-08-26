package entity

// CommentLike 评论点赞（与文章点赞 likes 对称）
type CommentLike struct {
	CreatedOnlyEntity
	CommentID int64 `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_comment_likes,priority:1;comment:评论ID" json:"comment_id"`
	UserID    int64 `gorm:"type:bigint unsigned;not null;uniqueIndex:uk_comment_likes,priority:2;comment:点赞用户ID" json:"user_id"`
}

// TableName 指定表名
func (CommentLike) TableName() string {
	return "comment_likes"
}
