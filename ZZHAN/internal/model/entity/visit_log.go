package entity

// VisitLog 访问日志（PV/UV 统计来源）
type VisitLog struct {
	CreatedOnlyEntity
	Path       string `gorm:"type:varchar(255);index;comment:访问路径" json:"path"`
	IP         string `gorm:"type:varchar(45);comment:客户端 IP" json:"ip"`
	VisitorKey string `gorm:"type:varchar(64);comment:访客指纹(UV 去重)" json:"visitor_key"`
	UserAgent  string `gorm:"type:varchar(255);comment:浏览器 UA" json:"user_agent"`
}

// TableName 指定表名
func (VisitLog) TableName() string {
	return "visit_logs"
}
