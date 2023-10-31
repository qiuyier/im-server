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

func init() {
	ack = &Ack{}
	ack.timeWheel = time_wheel.NewTimeWheel[*AckContent](ack._handler)
}

func (a *Ack) Start(ctx context.Context) error {
	_ = grpool.AddWithRecover(ctx, func(ctx context.Context) {
		a.timeWheel.Start()
	}, nil)

	<-ctx.Done()

	a.timeWheel.Stop()

	return gerror.New("ack service stopped")
}

func (a *Ack) insert(ackKey string, value *AckContent) {
	a.timeWheel.AddOnce(ackKey, 4*time.Second, value)
}

func (a *Ack) delete(ackKey string) {
	a.timeWheel.Remove(ackKey)
}

func (a *Ack) _handler(_ *time_wheel.DefaultTimeWheel[*AckContent], _ string, ackContent *AckContent) {
	channel, ok := Session.Channel(ackContent.channel)
	if !ok {
		return
	}

	client, ok := channel.Client(ackContent.cid)
	if !ok {
		return
	}

	if client.Closed() || gconv.Int64(client.uid) != ackContent.uid {
		return
	}

	if err := client.Write(ackContent.response); err != nil {
		g.Log().Errorf(gctx.New(), "ack err: ", err)
	}
}
