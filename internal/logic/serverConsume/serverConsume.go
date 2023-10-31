// Package serverConsume 处理服务端下发消息，订阅消息的消费逻辑实现
package serverConsume

import (
	"context"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"im/internal/consts"
	"im/internal/model"
	"im/internal/service"
	"im/utility/cache"
	"im/utility/websocket"
)

type sServerConsume struct {
	handlers      map[string]func(ctx context.Context, data []byte)
	clientStorage *cache.ClientCache
}

func init() {
	service.RegisterServerConsume(New())
}

func New() *sServerConsume {
	return &sServerConsume{
		clientStorage: cache.NewClientCache(),
	}
}

func (s *sServerConsume) Call(ctx context.Context, event string, data []byte) {
	if s.handlers == nil {
		s.init()
	}

	if call, ok := s.handlers[event]; ok {
		call(ctx, data)
	} else {
		g.Log().Errorf(ctx, "consume chat event: [%s]未注册回调事件", event)
	}
}

func (s *sServerConsume) init() {
	s.handlers = make(map[string]func(ctx context.Context, data []byte))

	s.handlers[consts.SubEventImMessageFoo] = s.onPublishFoo
	s.handlers[consts.SubEventImMessageBar] = s.onPublishBar
}

// 处理 serverEvent.onFoo 发布的消息
func (s *sServerConsume) onPublishFoo(ctx context.Context, data []byte) {
	if j, err := gjson.DecodeToJson(data); err != nil {
		g.Log().Error(ctx, "[Subscribe] onPublishFoo Unmarshal err: ", err.Error())
	} else {
		var in model.FooMessage
		if err := j.Scan(&in); err != nil {
			g.Log().Error(ctx, "[Subscribe] onPublishFoo Unmarshal err: ", err.Error())
		}

		clientIds := s.clientStorage.GetClientIdsFromUid(ctx, websocket.Session.Foo.Name(), gconv.String(in.ReceiverId))
		if len(clientIds) == 0 {
			return
		}

		content := websocket.NewSenderContent()
		content.SetAck(false)
		content.SetReceiver(clientIds...)
		content.SetMessage(consts.PushEventImMessageFoo, g.Map{
			"foo": "foo",
		})

		websocket.Session.Foo.Write(content)
	}
}

func (s *sServerConsume) onPublishBar(ctx context.Context, data []byte) {
	if j, err := gjson.DecodeToJson(data); err != nil {
		g.Log().Error(ctx, "[Subscribe] onPublishBar Unmarshal err: ", err.Error())
	} else {
		var in model.BarMessage
		if err := j.Scan(&in); err != nil {
			g.Log().Error(ctx, "[Subscribe] onPublishBar Unmarshal err: ", err.Error())
		}

		clientIds := s.clientStorage.GetClientIdsFromUid(ctx, websocket.Session.Foo.Name(), gconv.String(in.ReceiverId))
		if len(clientIds) == 0 {
			return
		}

		content := websocket.NewSenderContent()
		content.SetAck(true)
		content.SetReceiver(clientIds...)
		content.SetMessage(consts.PushEventImMessageBar, g.Map{
			"bar":       "bar",
			"record_id": in.RecordId,
			"bar_type":  in.BarType,
		})

		websocket.Session.Foo.Write(content)
	}
}
