package entity

// AboutItem 关于我-统一条目（技术栈 skills / 动态 dynamics；projects/experiences/roadmap 已移除）
// content 为 JSON 扩展字段，按 section 不同取值不同：
//   - skills:   {"level": 92}
//   - dynamics: {"text": "...", "link": "..."}
type AboutItem struct {
	CreatedOnlyEntity
	Section   string  `gorm:"type:varchar(20);not null;index:idx_about_section,priority:1;comment:分区: skills/dynamics" json:"section"`
	Type      *string `gorm:"type:varchar(20);comment:子类型: dynamics=write/star/talk; skills 为 NULL" json:"type"`
	Title     *string `gorm:"type:varchar(150);comment:名称/标题(skills.name / dynamics 无)" json:"title"`
	Content   string  `gorm:"type:json;comment:扩展字段(level/text/link 等,按 section 不同)" json:"content"`
	SortOrder int     `gorm:"type:int;not null;default:0;index:idx_about_section,priority:2;comment:排序值(同 section 内,越小越靠前)" json:"sort_order"`
}

// TableName 指定表名
func (AboutItem) TableName() string {
	return "about_items"
}
