package dto

import "time"

// CommentItem 评论列表项
type CommentItem struct {
	ID         int64         `json:"id"`
	ParentID   *int64        `json:"parent_id"`
	UserName   string        `json:"user_name"`
	UserAvatar string        `json:"avatar"`
	Content    string        `json:"content"`
	CreatedAt  time.Time     `json:"time"`
	LikeCount  int32         `json:"like_count"`
	Liked      bool          `json:"liked"`
	Replies    []CommentItem `json:"replies,omitempty"`
}

// CommentCreateResult 发表评论返回
type CommentCreateResult struct {
	ID      int64  `json:"id"`
	Status  string `json:"status"`
	Message string `json:"message"`
}

// CommentListResponse 评论列表响应
type CommentListResponse struct {
	List  []CommentItem `json:"list"`
	Total int64         `json:"total"`
}
