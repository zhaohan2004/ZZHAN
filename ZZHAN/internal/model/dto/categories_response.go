package dto

// CategoryItem 分类列表项
type CategoryItem struct {
	ID           int64  `json:"id"`    // 分类ID
	Name         string `json:"name"`  // 分类名称
	Slug         string `json:"slug"`  // URL 别名
	Icon         string `json:"icon"`  // 图标名
	Description  string `json:"desc"`  // 分类描述
	Color        string `json:"color"` // 主题色
	ArticleCount int64  `json:"count"` // 文章数量
}
