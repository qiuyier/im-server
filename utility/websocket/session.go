// Package websocket 用于管理 WebSocket 连接的会话 负责初始化会话对象、渠道和一些守护协程
package websocket

import (
	"context"
	"golang.org/x/sync/errgroup"
	"sync"
	"time"
)

var (
	Session *session
	once    sync.Once
)

type session struct {
	Foo      *Channel
	channels map[string]*Channel
}

func Init(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	once.Do(func() {
		initialize(ctx, eg, fn)
	})
}

func initialize(ctx context.Context, eg *errgroup.Group, fn func(name string)) {
	Session = &session{
		Foo:      NewChannel("foo", make(chan *SenderContent, 5<<20)), // 创建了一个带缓冲区的通道, 大小为5*2^20
		channels: map[string]*Channel{},
	}

	Session.channels["foo"] = Session.Foo

	// 延时启动
	time.AfterFunc(3*time.Second, func() {
		// 启动心跳监测
		eg.Go(func() error {
			defer fn("health exit")
			return health.Start(ctx)
		})

		// 启动应答机制
		eg.Go(func() error {
			defer fn("ack exit")
			return ack.Start(ctx)
		})

		// 启动渠道消费
		eg.Go(func() error {
			defer fn("channel consume exit")
			return Session.Foo.Start(ctx)
		})
	})
}

func (s *session) Channel(name string) (*Channel, bool) {
	val, ok := s.channels[name]
	return val, ok
}
