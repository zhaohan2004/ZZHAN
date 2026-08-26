package dto

// SiteResponse 站点响应
type SiteResponse struct {
	Name         string        `json:"name"`
	LogoText     string        `json:"logo_text"`
	Tagline      string        `json:"tagline"`
	Bio          string        `json:"bio"`
	Github       string        `json:"github"`
	Email        string        `json:"email"`
	Socials      []SocialsList `json:"socials"`
	Author       string        `json:"author"`
	Role         string        `json:"role"`
	Motto        string        `json:"motto"`
	Location     string        `json:"location"`
	Since        int           `json:"since"`
	Avatar       string        `json:"avatar"`
	HeroTerminal string        `json:"hero_terminal"`
}

// SocialsList 社交列表
type SocialsList struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
	Url  string `json:"url"`
}
