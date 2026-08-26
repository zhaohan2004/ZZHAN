package entity

import "time"

// 基础实体类型：依据 docs/database.sql 的列定义，
// 各表时间戳字段不尽相同，故拆分为四类基类按需嵌入。

// BaseEntity 标准实体：主键 + 创建时间 + 更新时间
type BaseEntity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键ID" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime;not null;comment:创建时间" json:"created_at"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;comment:更新时间" json:"updated_at"`
}

// CreatedOnlyEntity 仅含主键 + 创建时间
type CreatedOnlyEntity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键ID" json:"id"`
	CreatedAt time.Time `gorm:"type:datetime;not null;comment:创建时间" json:"created_at"`
}

// UpdatedOnlyEntity 仅含主键 + 更新时间
type UpdatedOnlyEntity struct {
	ID        int64     `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键ID" json:"id"`
	UpdatedAt time.Time `gorm:"type:datetime;not null;comment:更新时间" json:"updated_at"`
}

// IDOnlyEntity 仅含主键
type IDOnlyEntity struct {
	ID int64 `gorm:"primaryKey;autoIncrement;type:bigint unsigned;comment:主键ID" json:"id"`
}
