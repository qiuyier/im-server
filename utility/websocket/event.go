// Package websocket 定义一个事件处理框架，用于处理 WebSocket 连接的各种事件，如连接建立成功、接收消息、连接关闭和销毁连接
package websocket

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

type IEvent interface {
	// Open 客户端连接回调事件调用
	Open(client IClient)
	// Message 客户端消息回调事件调用
	Message(client IClient, data []byte)
	// Close 客户端连接关闭回调事件调用
	Close(client IClient, code int, text string)
}

type (
	OpenEvent    func(client IClient)
	MessageEvent func(client IClient, data []byte)
	CloseEvent   func(client IClient, code int, text string)
	EventOption  func(event *Event)
)

type Event struct {
	open    OpenEvent
	message MessageEvent
	close   CloseEvent
}

func NewEvent(opts ...EventOption) IEvent {
	o := &Event{}

	// 绑定自定义回调事件
	for _, opt := range opts {
		opt(o)
	}

	return o
}

func (e *Event) Open(client IClient) {
	if e.open == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			g.Log().Errorf(gctx.New(), "open event callback exception: ", client.Uid(), client.Cid(), client.Channel().Name(), err)
		}
	}()

	e.open(client)
}

func (e *Event) Message(client IClient, data []byte) {
	if e.message == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			g.Log().Errorf(gctx.New(), "message event callback exception: ", client.Uid(), client.Cid(), client.Channel().Name(), err)
		}
	}()

	e.message(client, data)
}

func (e *Event) Close(client IClient, code int, text string) {
	if e.close == nil {
		return
	}

	defer func() {
		if err := recover(); err != nil {
			g.Log().Errorf(gctx.New(), "close event callback exception: ", client.Uid(), client.Cid(), client.Channel().Name(), err)
		}
	}()

	e.close(client, code, text)
}

// WithOpenEvent 绑定open事件
func WithOpenEvent(e OpenEvent) EventOption {
	return func(event *Event) {
		event.open = e
	}
}

// WithMessageEvent 绑定message事件
func WithMessageEvent(e MessageEvent) EventOption {
	return func(event *Event) {
		event.message = e
	}
}

// WithCloseEvent 绑定close事件
func WithCloseEvent(e CloseEvent) EventOption {
	return func(event *Event) {
		event.close = e
	}
}
