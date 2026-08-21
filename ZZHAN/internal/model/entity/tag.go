package entity

// Tag 标签
type Tag struct {
	CreatedOnlyEntity
	Name string `gorm:"type:varchar(50);not null;uniqueIndex;comment:标签名称" json:"name"`
	Slug string `gorm:"type:varchar(60);not null;uniqueIndex;comment:URL 别名" json:"slug"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
