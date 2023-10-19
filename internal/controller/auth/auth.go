package auth

import (
	"context"
	"im/api/auth"
	"im/internal/model"
	"im/internal/service"
)

type cAuth struct {
}

var CAuth = cAuth{}

func (c *cAuth) Login(ctx context.Context, req *auth.LoginReq) (res *auth.LoginRes, err error) {
	out, err := service.Auth().Login(ctx, model.LoginInput{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		return nil, err
	}

	return &auth.LoginRes{Res: out}, nil
}
