package entity

// Admin 后台管理员表（单人，仅一条记录；账号/密码可由「个人信息」修改）
type Admin struct {
	BaseEntity
	Username     string `gorm:"type:varchar(50);not null;uniqueIndex;comment:后台登录名（可修改）" json:"username"`
	PasswordHash string `gorm:"type:varchar(255);not null;comment:密码哈希(bcrypt，可修改)" json:"-"`
	Nickname     string `gorm:"type:varchar(50);comment:显示昵称" json:"nickname"`
	Avatar       string `gorm:"type:text;comment:头像 URL / base64" json:"avatar"`
}

// TableName 指定表名
func (Admin) TableName() string {
	return "admins"
}
