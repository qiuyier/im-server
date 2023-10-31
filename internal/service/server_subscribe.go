// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"im/utility/websocket"
	"net/http"

	"golang.org/x/sync/errgroup"
)

type (
	IServerSubscribe interface {
		Conn(w http.ResponseWriter, r *http.Request) error
		NewClient(uid int, conn websocket.ISocket) error
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
