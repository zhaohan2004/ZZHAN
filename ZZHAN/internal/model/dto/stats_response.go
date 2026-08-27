package dto

// Dynamic 动态项
type Dynamic struct {
	Type string `json:"type"`
	Text string `json:"text"`
	Time string `json:"time"`
	Link string `json:"link,omitempty"`
}

// StatsResponse 站点统计响应
type StatsResponse struct {
	Articles int       `json:"articles"`
	Views    int64     `json:"views"`
	Comments int64     `json:"comments"`
	Dynamics []Dynamic `json:"dynamics"`
}
