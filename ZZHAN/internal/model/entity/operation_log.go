package entity

// OperationLog 后台操作记录
type OperationLog struct {
	CreatedOnlyEntity
	AdminID *uint64 `gorm:"type:bigint unsigned;comment:操作者(后台管理员ID)" json:"adminId"`
	Action  string  `gorm:"type:varchar(50);not null;comment:操作类型(create/update/delete 等)" json:"action"`
	Target  string  `gorm:"type:varchar(255);comment:操作对象(如文章标题/ID)" json:"target"`
}

// TableName 指定表名
func (OperationLog) TableName() string {
	return "operation_logs"
}
