// Package serverEvent 自定义回调事件具体实现
package serverEvent

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"im/internal/consts"
	"im/internal/model"
	"im/internal/service"
	"im/utility/websocket"
)

type sServerEvent struct {
	handlers map[string]func(ctx context.Context, client websocket.IClient, data []byte)
}

func init() {
	service.RegisterServerEvent(New())
}

func New() *sServerEvent {
	return &sServerEvent{}
}

// 初始化客户端自定义事件
func (s *sServerEvent) init() {
	s.handlers = make(map[string]func(ctx context.Context, client websocket.IClient, data []byte))

	// 注册自定义事件
	s.handlers["im.message.foo"] = s.onFoo
}

// Call 匹配处理客户端请求事件
func (s *sServerEvent) Call(ctx context.Context, client websocket.IClient, event string, data []byte) {
	if s.handlers == nil {
		s.init()
	}

	if call, ok := s.handlers[event]; ok {
		call(ctx, client, data)
	} else {
		g.Log().Errorf(ctx, "Event: [%s]未注册回调事件", event)
	}
}

// onFoo 示例
func (s *sServerEvent) onFoo(ctx context.Context, client websocket.IClient, data []byte) {
	_, err := g.Redis().Publish(ctx, consts.ImTopicBar, gjson.MustEncodeString(&model.SubscribeContent{
		Event: consts.SubEventImMessageFoo,
		Data: gjson.MustEncodeString(model.FooMessage{
			ReceiverId: client.Uid(),
			AcceptData: gconv.String(data),
		}),
	}))
	if err != nil {
		g.Log().Errorf(ctx, "onFoo publish err: ", err)
		return
	}
}

// OnOpen 客户端连接回调事件
func (s *sServerEvent) OnOpen(client websocket.IClient) {
	g.Log().Infof(gctx.New(), "%d goes online", client.Cid())
}

// OnMessage 客户端消息处理回调事件
func (s *sServerEvent) OnMessage(client websocket.IClient, data []byte) {
	j, err := gjson.LoadJson(data)
	if err != nil {
		g.Log().Errorf(gctx.New(), "onMessage json unmarshall err: ", err)
		return
	}

	if !j.Get("event").IsEmpty() {
		s.Call(gctx.New(), client, j.Get("event").String(), data)
	}
}

// OnClose 客户端关闭回调事件
func (s *sServerEvent) OnClose(client websocket.IClient, code int, text string) {
	g.Log().Infof(gctx.New(), "client close uid: %d, cid: %d, channel: %s, code: %d, text: %s", client.Uid(), client.Cid(), client.Channel().Name(), code, text)
}
