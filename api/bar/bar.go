// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package bar

import (
	"context"

	"im/api/bar/v1"
)

type IBarV1 interface {
	Bar(ctx context.Context, req *v1.BarReq) (res *v1.BarRes, err error)
}
