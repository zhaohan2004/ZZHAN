package entity

// Category 文章分类
type Category struct {
	BaseEntity
	Name        string `gorm:"type:varchar(50);not null;index;comment:分类名称" json:"name"`
	Slug        string `gorm:"type:varchar(60);not null;uniqueIndex;comment:URL 别名" json:"slug"`
	Icon        string `gorm:"type:varchar(30);comment:lucide 图标名" json:"icon"`
	Description string `gorm:"type:varchar(255);comment:分类描述" json:"description"`
	Color       string `gorm:"type:varchar(10);default:#3b82f6;comment:主题色(十六进制)" json:"color"`
	SortOrder   int    `gorm:"type:int;not null;default:0;comment:排序值(越小越靠前)" json:"sort_order"`
	Status      string `gorm:"type:varchar(20);not null;default:active;comment:状态(active启用/inactive禁用)" json:"status"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}
