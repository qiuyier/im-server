// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"im/internal/model"
)

type (
	IBar interface {
		Bar(ctx context.Context, in model.BarReq) (out *model.BarRes, err error)
	}
)

var (
	localBar IBar
)

func Bar() IBar {
	if localBar == nil {
		panic("implement not found for interface IBar, forgot register?")
	}
	return localBar
}

func RegisterBar(i IBar) {
	localBar = i
}
