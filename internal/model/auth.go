package model

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	// Token 类型
	Type string `json:"type,omitempty"`
	// token
	AccessToken string `json:"access_token,omitempty"`
	// 过期时间
	ExpiresIn int `json:"expires_in,omitempty"`
}
