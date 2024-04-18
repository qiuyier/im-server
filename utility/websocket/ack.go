// Package websocket 消息确认机制，当客户端及时响应后把任务从时间轮删除，反之则重复下发
package websocket

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gconv"
	"im/utility/time_wheel"
	"time"
)

type AckContent struct {
	cid      int64
	uid      int64
	channel  string
	response *ClientResponse
}

type Ack struct {
	timeWheel *time_wheel.DefaultTimeWheel[*AckContent]
}

var ack *Ack

// 初始化时间轮
func init() {
	ack = &Ack{}
	ack.timeWheel = time_wheel.NewTimeWheel[*AckContent](ack._handler)
}

// Start 启动消息确认服务
func (a *Ack) Start(ctx context.Context) error {
	_ = grpool.AddWithRecover(ctx, func(ctx context.Context) {
		a.timeWheel.Start()
	}, nil)

	<-ctx.Done()

	a.timeWheel.Stop()

	return gerror.New("ack service stopped")
}

// 添加消息确认事件
func (a *Ack) insert(ackKey string, value *AckContent) {
	a.timeWheel.AddOnce(ackKey, 10*time.Second, value)
}

// 删除消息确认事件
func (a *Ack) delete(ackKey string) {
	a.timeWheel.Remove(ackKey)
}

// 消息确认事件实现逻辑
func (a *Ack) _handler(_ *time_wheel.DefaultTimeWheel[*AckContent], _ string, ackContent *AckContent) {
	// 判断消息渠道是否存在
	channel, ok := Session.Channel(ackContent.channel)
	if !ok {
		return
	}

	// 获取渠道内用户客户端
	client, ok := channel.Client(ackContent.cid)
	if !ok {
		return
	}

	// 判断收信人是否一致，判断客户端是否关闭
	if client.Closed() || gconv.Int64(client.uid) != ackContent.uid {
		return
	}

	// 把需要下发的消息添加到发送通道，等待发送
	if err := client.Write(ackContent.response); err != nil {
		g.Log().Errorf(gctx.New(), "ack handler err: %v", err)
	}
}
