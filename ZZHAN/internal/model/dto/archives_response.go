package dto

// ArchiveArticle 归档文章项
type ArchiveArticle struct {
	ID       int64  `json:"id"`
	Slug     string `json:"slug"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Category string `json:"category"`
	Views    int32  `json:"views"`
}

// ArchiveItem 归档月份项
type ArchiveItem struct {
	Year     string           `json:"year"`
	Month    string           `json:"month"`
	Count    int              `json:"count"`
	Articles []ArchiveArticle `json:"articles"`
}
