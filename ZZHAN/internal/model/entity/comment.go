package entity

// Comment 评论表
type Comment struct {
	CreatedOnlyEntity
	ArticleID  uint64  `gorm:"type:bigint unsigned;not null;index:idx_comments_article,priority:1;comment:所属文章ID" json:"article_id"`
	ParentID   *uint64 `gorm:"type:bigint unsigned;index;comment:父评论ID(楼中楼回复)" json:"parent_id"`
	UserID     *uint64 `gorm:"type:bigint unsigned;comment:登录用户ID；匿名为 NULL" json:"user_id"`
	UserName   string  `gorm:"type:varchar(50);comment:匿名访客昵称" json:"user_name"`
	UserAvatar string  `gorm:"type:varchar(255);comment:访客头像 URL" json:"user_avatar"`
	UserEmail  string  `gorm:"type:varchar(100);comment:访客邮箱" json:"user_email"`
	Content    string  `gorm:"type:text;not null;comment:评论内容" json:"content"`
	IP         string  `gorm:"type:varchar(45);comment:评论者 IP" json:"ip"`
	Status     string  `gorm:"type:varchar(20);not null;default:normal;index:idx_comments_article,priority:2;comment:状态：normal(正常)/banned(封禁)" json:"status"`
	LikeCount  uint    `gorm:"type:int unsigned;not null;default:0;comment:点赞数(冗余计数)" json:"like_count"`
}

// TableName 指定表名
func (Comment) TableName() string {
	return "comments"
}
