package bar

import (
	"context"
	"im/internal/service"

	"im/api/bar/v1"
)

func (c *ControllerV1) Bar(ctx context.Context, req *v1.BarReq) (res *v1.BarRes, err error) {
	out, err := service.Bar().Bar(ctx, req.BarReq)
	if err != nil {
		return nil, err
	}

	res = &v1.BarRes{
		BarRes: out,
	}

	return res, nil
}
