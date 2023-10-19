package auth

import "github.com/gogf/gf/v2/frame/g"

type LoginReq struct {
	g.Meta   `path:"/login" method:"post" summary:"登录" tags:"Login"`
	Username string `json:"username" v:"required" dc:"用户名"`
	Password string `json:"password" v:"required" dc:"密码"`
}

type LoginRes struct {
	Res any `json:"res"`
}
