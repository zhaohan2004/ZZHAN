package dto

// CommentsAdminItem 后台评论列表项
type CommentsAdminItem struct {
	ID           int64  `json:"id"`
	ArticleID    int64  `json:"article_id"`
	ArticleTitle string `json:"article_title"` // 所属文章标题
	ParentID     *int64 `json:"parent_id"`
	UserID       int64  `json:"user_id"`
	UserName     string `json:"user_name"`
	UserAvatar   string `json:"user_avatar"`
	Content      string `json:"content"`
	Status       string `json:"status"`
	LikeCount    int32  `json:"like_count"`
	IP           string `json:"ip"`
	CreatedAt    string `json:"created_at"`
}

// CommentsAdminListResponse 后台评论列表响应
type CommentsAdminListResponse struct {
	List  []CommentsAdminItem `json:"list"`
	Total int64               `json:"total"`
}
