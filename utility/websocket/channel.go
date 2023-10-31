// Package websocket WebSocket 通道管理器，用于管理客户端连接和消息传递
package websocket

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/orcaman/concurrent-map/v2"
	"github.com/sourcegraph/conc/pool"
	"sync/atomic"
	"time"
)

type IChannel interface {
	// Name 频道名
	Name() string

	// Count 统计客户端连接数
	Count() int64

	// Client 获取客户端
	Client(cid int64) (*Client, bool)

	// Write 推送消息到消费通道
	Write(data *SenderContent)

	// addClient 添加客户端
	addClient(client *Client)

	// delClient 删除客户端
	delClient(client *Client)
}

// Channel 实现了 IChannel 接口
type Channel struct {
	name    string                              // 通道名称
	count   int64                               // 客户端连接数
	node    cmap.ConcurrentMap[string, *Client] // 客户端列表
	outChan chan *SenderContent                 // 消息发送通道
}

// NewChannel 创建一个新的通道实例，需要指定通道名称和消息发送通道
func NewChannel(name string, outChan chan *SenderContent) *Channel {
	return &Channel{
		name:    name,
		node:    cmap.New[*Client](),
		outChan: outChan,
	}
}

func (c *Channel) Name() string {
	return c.name
}

func (c *Channel) Count() int64 {
	return c.count
}

func (c *Channel) Client(cid int64) (*Client, bool) {
	return c.node.Get(gconv.String(cid))
}

func (c *Channel) Write(data *SenderContent) {
	timer := time.NewTimer(5 * time.Second)
	defer timer.Stop()

	select {
	case c.outChan <- data:
	case <-timer.C:
		g.Log().Errorf(gctx.New(), "[%s] Channel OutChan 写入消息超时, 管道长度: %d", c.name, len(c.outChan))
	}
}

// 添加客户端到通道
func (c *Channel) addClient(client *Client) {
	c.node.Set(gconv.String(client.cid), client)
	atomic.AddInt64(&c.count, 1)
}

// 从通道中删除客户端
func (c *Channel) delClient(client *Client) {
	cid := gconv.String(client.cid)

	if !c.node.Has(cid) {
		return
	}

	c.node.Remove(cid)

	atomic.AddInt64(&c.count, -1)
}

// Start 通道的主要处理逻辑
// 用于启动一个守护协程来处理消息发送和客户端的管理
// 使用定时器定期检查通道消息发送通道中的消息，以及处理消息的消费
func (c *Channel) Start(ctx context.Context) error {
	var (
		worker = pool.New().WithMaxGoroutines(10)
		timer  = time.NewTicker(15 * time.Second)
	)
	defer timer.Stop()
	defer g.Log().Errorf(ctx, "channel exit: %s", c.Name())

	for {
		select {
		case <-ctx.Done():
			// 通道管理器被关闭的信号，如果传入的上下文 ctx 被取消，说明通道需要退出，此时会返回一个错误以指示通道退出。
			return gerror.Newf("channel exit: %s", c.Name())
		case <-timer.C:
			//定时器触发，当没有消息时记录日志
			g.Log().Debugf(ctx, "channel empty message name: %s, len: %d", c.name, len(c.outChan))
		case value, ok := <-c.outChan:
			// 如果接收到 <-c.outChan，说明有消息需要发送。消息会交给 c.consume 来处理
			if !ok {
				return gerror.Newf("outChan close: %s", c.Name())
			}

			c.consume(worker, value, func(data *SenderContent, client *Client) {
				_ = client.Write(&ClientResponse{
					IsAck:   data.IsAck,
					Event:   data.message.Event,
					Content: data.message.Content,
					Retry:   3,
				})
			})
		}
	}
}

// 如果有消息需要发送给客户端，使用协程池并发地处理消息发送操作
// 如果是广播消息，它会将消息发送给通道中的所有客户端；如果是定向消息，它会根据指定的客户端唯一标识将消息发送给指定客户端
func (c *Channel) consume(worker *pool.Pool, data *SenderContent, fn func(data *SenderContent, client *Client)) {
	worker.Go(func() {
		if data.IsBroadcast() {
			// 遍历c.node
			c.node.IterCb(func(_ string, client *Client) {
				fn(data, client)
			})
			return
		}

		for _, cid := range data.receivers {
			if client, ok := c.Client(cid); ok {
				fn(data, client)
			}
		}
	})
}
