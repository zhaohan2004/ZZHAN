package dto

// TagListItem 标签列表项
type TagListItem struct {
	ID           int64  `json:"id"`    // 标签ID
	Name         string `json:"name"`  // 标签名称
	ArticleCount int64  `json:"count"` // 文章数量
}
