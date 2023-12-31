// Package serverSubscribe websocket服务一部分
// 主要负责处理连接建立、消息订阅以及消息的处理和分发。
// 通过消息订阅，可以接收来自 Redis 的消息，然后将其分发给连接的客户端，实现了 WebSocket 通信的功能
package serverSubscribe

import (
	"context"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gconv"
	"golang.org/x/sync/errgroup"
	"im/internal/consts"
	"im/internal/model"
	"im/internal/service"
	"im/utility/cache"
	"im/utility/websocket"
	"sync"
)

var once sync.Once

type sServerSubscribe struct {
	clientStorage *cache.ClientCache
}

func init() {
	service.RegisterServerSubscribe(New())
}

func New() *sServerSubscribe {
	return &sServerSubscribe{
		clientStorage: cache.NewClientCache(),
	}
}

// Conn 建立websocket连接
func (s *sServerSubscribe) Conn(r *ghttp.Request) error {
	// 将请求升级为websocket服务，并获取连接对象
	conn, err := websocket.NewGfWebSocket(r)
	if err != nil {
		g.Log().Errorf(r.Context(), "websocket connect err: ", err)
		return err
	}

	return s.NewClient(service.Session().GetUid(r.Request.Context()), conn)
}

// NewClient 将连接保存为自定义的客户端对象
func (s *sServerSubscribe) NewClient(uid int, conn websocket.ISocket) error {
	return websocket.NewClient(conn, &websocket.ClientOption{
		Uid:     uid,
		Channel: websocket.Session.Foo,
		Storage: s.clientStorage,
		Buffer:  10,
	}, websocket.NewEvent(
		// 注册连接成功回调事件
		websocket.WithOpenEvent(service.ServerEvent().OnOpen),

		// 注册接收消息回调事件
		websocket.WithMessageEvent(service.ServerEvent().OnMessage),

		// 注册关闭连接毁掉事件
		websocket.WithCloseEvent(service.ServerEvent().OnClose),
	))
}

// Start 启动服务监听
func (s *sServerSubscribe) Start(ctx context.Context, eg *errgroup.Group) {
	// 使用Once保证全局只有一个监听器
	// 不然假如有多个订阅者，redis会并行地将消息发送给所有监听者，导致同一消息会被重复消费
	once.Do(func() {
		eg.Go(func() error {
			return s.SetUpMessageSubscribe(ctx)
		})
	})
}

func (s *sServerSubscribe) SetUpMessageSubscribe(ctx context.Context) error {
	_ = grpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		s.subscribe(ctx, []string{consts.ImTopicBar}, service.ServerConsume())
	}, nil)

	<-ctx.Done()

	return nil
}

// redis订阅消息逻辑
func (s *sServerSubscribe) subscribe(ctx context.Context, topic []string, consume service.IServerConsume) {
	defaultTopic := consts.ImTopicFoo
	conn, err := g.Redis().Conn(ctx)
	if err != nil {
		g.Log().Errorf(ctx, "redis sub con err: ", err)
	}
	defer func() {
		err := conn.Close(ctx)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}
	}()

	_, err = conn.Subscribe(ctx, defaultTopic, topic...)
	if err != nil {
		g.Log().Errorf(ctx, "redis sub err: ", err)
	}

	g.Log().Infof(ctx, "Start MessageSubscribe...")

	for {
		msg, err := conn.ReceiveMessage(ctx)
		if err != nil {
			g.Log().Fatal(ctx, err)
		}
		s.handle(ctx, msg, consume)
	}
}

// 处理订阅消息逻辑
func (s *sServerSubscribe) handle(ctx context.Context, data *gredis.Message, consume service.IServerConsume) {
	_ = grpool.Add(ctx, func(ctx context.Context) {
		if j, err := gjson.DecodeToJson(data.Payload); err != nil {
			panic(err)
		} else {
			var in model.SubscribeContent
			if err := j.Scan(&in); err != nil {
				panic(err)
			}
			consume.Call(ctx, in.Event, gconv.Bytes(in.Data))
		}
		defer func() {
			if err := recover(); err != nil {
				g.Log().Error(ctx, "MessageSubscribe Call Err:", err)
			}
		}()

	})
}
