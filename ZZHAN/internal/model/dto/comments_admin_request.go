package dto

// CommentsAdminListRequest 后台评论列表请求
type CommentsAdminListRequest struct {
	Page      int    `form:"page" binding:"omitempty,min=1"`
	PageSize  int    `form:"page_size" binding:"omitempty,min=1,max=100"`
	Keyword   string `form:"keyword"`                                            // 搜索关键词（内容/用户名）
	Status    string `form:"status" binding:"omitempty,oneof=normal banned all"` // 按状态筛选
	ArticleID int64  `form:"article_id"`                                         // 按文章筛选
	StartDate string `form:"start_date"`                                         // 开始日期（2006-01-02）
	EndDate   string `form:"end_date"`                                           // 结束日期（2006-01-02）
}

// CommentsAdminUpdateStatusRequest 后台修改评论状态请求
type CommentsAdminUpdateStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=normal banned"`
}
