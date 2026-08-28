package entity

// Tag 标签
type Tag struct {
	BaseEntity
	Name   string `gorm:"type:varchar(50);not null;uniqueIndex;comment:标签名称" json:"name"`
	Slug   string `gorm:"type:varchar(60);not null;uniqueIndex;comment:URL 别名" json:"slug"`
	Status string `gorm:"type:varchar(20);not null;default:active;comment:状态(active启用/inactive禁用)" json:"status"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}
