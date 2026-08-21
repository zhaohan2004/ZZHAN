package entity

import "time"

// User 前台用户表（微信 / GitHub OAuth 登录，无密码）
type User struct {
	BaseEntity
	Provider    string     `gorm:"type:varchar(20);not null;uniqueIndex:uk_users_provider_openid,priority:1;comment:登录方式：wechat / github" json:"provider"`
	Openid      string     `gorm:"type:varchar(100);not null;uniqueIndex:uk_users_provider_openid,priority:2;comment:该平台用户唯一标识" json:"openid"`
	Unionid     string     `gorm:"type:varchar(100);index;comment:微信 unionid（跨应用唯一，GitHub 无）" json:"unionid"`
	Nickname    string     `gorm:"type:varchar(50);comment:昵称" json:"nickname"`
	Avatar      string     `gorm:"type:varchar(255);comment:头像 URL" json:"avatar"`
	Email       string     `gorm:"type:varchar(100);comment:邮箱" json:"email"`
	Status      int8       `gorm:"type:tinyint;not null;default:1;comment:状态：1 正常 / 0 禁用" json:"status"`
	LastLoginAt *time.Time `gorm:"type:datetime;comment:最后登录时间" json:"lastLoginAt"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}
