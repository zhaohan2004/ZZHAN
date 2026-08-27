package dto

// Skill 技能项
type Skill struct {
	Name  string `json:"name"`
	Level int    `json:"level"`
}

// AboutResponse 关于我响应
type AboutResponse struct {
	Skills []Skill `json:"skills"`
}
