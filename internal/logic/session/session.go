package session

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
	"im/internal/consts"
	"im/internal/service"
)

type sSession struct {
}

func init() {
	service.RegisterSession(New())
}

func New() *sSession {
	return &sSession{}
}

func (s *sSession) GetUid(ctx context.Context) int {
	return gconv.Int(ctx.Value(consts.CtxUserId))
}
