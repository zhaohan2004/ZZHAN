package dto

// StatItem 单项统计
type StatItem struct {
	Value int64 `json:"value"`
}

// RecentPostItem 最新文章
type RecentPostItem struct {
	ID       int64  `json:"id"`
	Title    string `json:"title"`
	Category string `json:"category"`
	Date     string `json:"date"`
	Views    int32  `json:"views"`
}

// DashboardCommentItem 最新评论
type DashboardCommentItem struct {
	ID       int64  `json:"id"`
	UserName string `json:"user_name"`
	Avatar   string `json:"avatar"`
	Content  string `json:"content"`
	Time     string `json:"time"`
}

// DashboardOperationItem 操作记录
type DashboardOperationItem struct {
	Time   string `json:"time"`
	User   string `json:"user"`
	Action string `json:"action"`
	Target string `json:"target"`
}
