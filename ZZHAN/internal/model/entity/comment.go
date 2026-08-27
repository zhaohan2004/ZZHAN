package entity

// Comment 评论表
type Comment struct {
	CreatedOnlyEntity
	ArticleID  int64  `gorm:"type:bigint unsigned;not null;index:idx_comments_article,priority:1;comment:所属文章ID" json:"article_id"`
	ParentID   *int64 `gorm:"type:bigint unsigned;index;comment:父评论ID" json:"parent_id"`
	UserID     int64  `gorm:"type:bigint unsigned;not null;comment:登录用户ID" json:"user_id"`
	UserName   string `gorm:"type:varchar(50);comment:用户昵称" json:"user_name"`
	UserAvatar string `gorm:"type:varchar(255);comment:用户头像 URL" json:"user_avatar"`
	Content    string `gorm:"type:text;not null;comment:评论内容" json:"content"`
	IP         string `gorm:"type:varchar(45);comment:评论者 IP" json:"ip"`
	Status     string `gorm:"type:varchar(20);not null;default:normal;index:idx_comments_article,priority:2;comment:状态：normal(正常)/banned(封禁)" json:"status"`
	LikeCount  int32  `gorm:"type:int unsigned;not null;default:0;comment:点赞数(冗余计数)" json:"like_count"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
