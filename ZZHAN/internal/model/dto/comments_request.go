package dto

// CommentListRequest 评论列表请求
type CommentListRequest struct {
	Page     int `form:"page" binding:"omitempty,min=1"`
	PageSize int `form:"page_size" binding:"omitempty,min=1,max=50"`
}

// CommentCreateRequest 发表评论请求
type CommentCreateRequest struct {
	Content  string `json:"content" binding:"required,min=1,max=1000"`
	ParentID *int64 `json:"parent_id" binding:"omitempty"`
}
