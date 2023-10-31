// ================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// You can delete these comments if you wish manually maintain this interface file.
// ================================================================================

package service

import (
	"context"
	"im/utility/websocket"
)

type (
	IServerEvent interface {
		Call(ctx context.Context, client websocket.IClient, event string, data []byte)
		OnOpen(client websocket.IClient)
		OnMessage(client websocket.IClient, data []byte)
		OnClose(client websocket.IClient, code int, text string)
	}
)

var (
	localServerEvent IServerEvent
)

func ServerEvent() IServerEvent {
	if localServerEvent == nil {
		panic("implement not found for interface IServerEvent, forgot register?")
	}
	return localServerEvent
}

func RegisterServerEvent(i IServerEvent) {
	localServerEvent = i
}
