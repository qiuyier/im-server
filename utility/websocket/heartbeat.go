// Package websocket 这段代码实现了一个客户端心跳管理模块，用于定期检测客户端的心跳，确保客户端的连接处于活跃状态
package websocket

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gconv"
	"im/utility/time_wheel"
	"time"
)

const (
	heartbeatInterval = 30 // 心跳检测间隔时间
	heartbeatTimeout  = 75 // 心跳检测超时时间(超时时间是隔间检测时间的2.5倍以上)
)

type heartbeat struct {
	timeWheel *time_wheel.DefaultTimeWheel[*Client]
}

var health *heartbeat

// 初始化时间轮
func init() {
	health = &heartbeat{}
	health.timeWheel = time_wheel.NewTimeWheel[*Client](health._handler)
}

// Start 启动客户端心跳管理监控
func (h *heartbeat) Start(ctx context.Context) error {
	_ = grpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		h.timeWheel.Start()
	}, nil)

	<-ctx.Done()

	h.timeWheel.Stop()

	return gerror.New("heartbeat exit")
}

// 向心跳时间轮中添加客户端
func (h *heartbeat) insert(c *Client) {
	h.timeWheel.AddOnce(gconv.String(c.cid), time.Duration(heartbeatInterval)*time.Second, c)
}

// 从心跳时间轮中移除客户端
func (h *heartbeat) delete(c *Client) {
	h.timeWheel.Remove(gconv.String(c.cid))
}

// 心跳时间轮的任务处理函数，用于处理心跳检测任务
func (h *heartbeat) _handler(timeWheel *time_wheel.DefaultTimeWheel[*Client], key string, c *Client) {
	// 检查客户端是否已关闭，如果已关闭则不进行处理
	if c.Closed() {
		return
	}

	// 计算客户端最后心跳时间与当前时间的时间间隔 interval，如果该间隔大于 heartbeatTimeout，表示心跳检测超时，会关闭客户端连接
	interval := int(time.Now().Unix() - c.lastTime)
	if interval > heartbeatTimeout {
		c.Close(20000, "心跳检测超时，自动关闭")
		return
	}

	// 超过心跳间隔时间则主动推送一次消息
	if interval >= heartbeatInterval {
		_ = c.Write(&ClientResponse{Event: "ping"})
	}

	timeWheel.AddOnce(key, time.Duration(heartbeatInterval)*time.Second, c)
}
