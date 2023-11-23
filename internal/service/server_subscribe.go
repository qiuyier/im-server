// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"im/utility/websocket"

	"github.com/gogf/gf/v2/net/ghttp"
	"golang.org/x/sync/errgroup"
)

type (
	IServerSubscribe interface {
		// Conn 建立websocket连接
		Conn(r *ghttp.Request) error
		// NewClient 将连接保存为自定义的客户端对象
		NewClient(uid int, conn websocket.ISocket) error
		// Start 启动服务监听
		Start(ctx context.Context, eg *errgroup.Group)
		SetUpMessageSubscribe(ctx context.Context) error
	}
)

var (
	localServerSubscribe IServerSubscribe
)

func ServerSubscribe() IServerSubscribe {
	if localServerSubscribe == nil {
		panic("implement not found for interface IServerSubscribe, forgot register?")
	}
	return localServerSubscribe
}

func RegisterServerSubscribe(i IServerSubscribe) {
	localServerSubscribe = i
}
