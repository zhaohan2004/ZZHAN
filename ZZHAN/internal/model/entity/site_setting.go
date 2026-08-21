package entity

// SiteSetting 站点设置（键值对）
type SiteSetting struct {
	UpdatedOnlyEntity
	Key         string `gorm:"type:varchar(80);not null;uniqueIndex;column:key;comment:设置键名" json:"key"`
	Value       string `gorm:"type:text;comment:设置值(JSON 或标量)" json:"value"`
	Description string `gorm:"type:varchar(255);comment:设置项说明" json:"description"`
}

// TableName 指定表名
func (SiteSetting) TableName() string {
	return "site_settings"
}
