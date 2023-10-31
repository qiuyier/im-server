package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"im/internal/model"
)

type BarReq struct {
	g.Meta `path:"/bar" tags:"Bar" method:"post" summary:"Bar"`
	model.BarReq
}
type BarRes struct {
	g.Meta `mime:"application/json" example:"json"`
	*model.BarRes
}
