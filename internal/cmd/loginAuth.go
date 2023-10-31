package cmd

import (
	"github.com/goflyfox/gtoken/gtoken"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/grand"
	"im/internal/config"
	"im/internal/consts"
	"im/utility/util"
)

type UserInfo struct {
	Username string `json:"username"`
	Password string `json:"password"`
	UserId   int    `json:"user_id"`
}

type AuthLoginRes struct {
	// Token 类型
	Type string `json:"type,omitempty"`
	// token
	AccessToken string `json:"access_token,omitempty"`
	// 过期时间
	ExpiresIn int `json:"expires_in,omitempty"`
}

func StartUserGToken() (userGToken *gtoken.GfToken, err error) {
	userGToken = &gtoken.GfToken{
		CacheMode:        2,
		ServerName:       "im-server",
		LoginPath:        "/login",
		LoginBeforeFunc:  loginFunc,
		LoginAfterFunc:   loginAfterFunc,
		LogoutPath:       "/logout",
		LogoutAfterFunc:  LogOutAfterFunc,
		AuthPaths:        g.SliceStr{},
		AuthExcludePaths: g.SliceStr{},
		AuthAfterFunc:    authAfterFunc,
		MultiLogin:       true,
	}
	return
}

func loginFunc(r *ghttp.Request) (string, interface{}) {
	username := r.Get("username").String()
	password := r.Get("password").String()

	if util.IsBlank(username) || util.IsBlank(password) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeMissingParameter.Code(),
			Message: "username or password is missing",
			Data:    nil,
		})
		r.ExitAll()
	}

	// 验证用户名密码
	if !util.ValidatePassword(username, password) {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeBusinessValidationFailed.Code(),
			Message: "username or password is not correct",
			Data:    nil,
		})
		r.ExitAll()
	}

	userinfo := &UserInfo{
		Username: username,
		Password: password,
		UserId:   grand.N(1, 10000),
	}

	return consts.GTokenUserPrefix + username, userinfo
}

func loginAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	if !respData.Success() {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeInternalError.Code(),
			Message: respData.Msg,
			Data:    nil,
		})
		r.ExitAll()
	} else {
		// 获得用户ID
		//userKey := respData.GetString("userKey")
		//merchantID := gstr.StrEx(userKey, consts.GTokenUserPrefix)

		data := &AuthLoginRes{
			Type:        "Bearer",
			AccessToken: respData.GetString("token"),
			ExpiresIn:   config.Cfg.Jwt.ExpiresTime,
		}

		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeOK.Code(),
			Message: "auth success",
			Data:    data,
		})
		r.ExitAll()
	}
	return
}

func LogOutAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	if !respData.Success() {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeInternalError.Code(),
			Message: respData.Msg,
			Data:    nil,
		})
		r.ExitAll()
	} else {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeOK.Code(),
			Message: "logout successful",
			Data:    nil,
		})
		r.ExitAll()
	}
	return
}

func authAfterFunc(r *ghttp.Request, respData gtoken.Resp) {
	var userInfo UserInfo
	err := gconv.Struct(respData.GetString("data"), &userInfo)
	if err != nil {
		r.Response.WriteJson(ghttp.DefaultHandlerResponse{
			Code:    gcode.CodeInternalError.Code(),
			Message: respData.Msg,
			Data:    nil,
		})
		r.ExitAll()
	}
	r.SetCtxVar(consts.CtxUserName, userInfo.Username)
	r.SetCtxVar(consts.CtxUserId, userInfo.UserId)

	r.Middleware.Next()
}
