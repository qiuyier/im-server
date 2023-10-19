package auth

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"im/internal/config"
	"im/internal/model"
	"im/internal/service"
	"im/utility/jwt"
	"im/utility/util"
)

type sAuth struct {
}

func init() {
	service.RegisterAuth(New())
}

func New() *sAuth {
	return &sAuth{}
}

func (s *sAuth) Login(ctx context.Context, in model.LoginInput) (out model.LoginOutput, err error) {
	if !util.ValidatePassword(in.Username, in.Password) {
		return model.LoginOutput{}, gerror.New("username or password invalid")
	}

	out = model.LoginOutput{
		Type:        "Bearer",
		AccessToken: jwt.GenerateJwtToken(in.Username, "user"),
		ExpiresIn:   config.Cfg.Jwt.ExpiresTime,
	}

	return
}
