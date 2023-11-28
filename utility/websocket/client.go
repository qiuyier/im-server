// Package websocket 实现 WebSocket 客户端的核心逻辑，包括消息的接收和发送、心跳检测、连接关闭等功能，以及回调方法的触发
package websocket

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"im/internal/consts"
	"im/utility/cache"
	"sync/atomic"
	"time"
)

// IClient 定义了客户端的相关操作
type IClient interface {
	Write(data *ClientResponse) error
	Close(code int, text string)
	Cid() int64        // 客户端ID
	Uid() int          // 客户端关联用户ID
	Channel() IChannel // 获取客户端所属渠道
}

type Client struct {
	conn     ISocket              // 客户端连接
	cid      int64                // 客户端ID/客户端唯一标识
	uid      int                  // 用户ID
	lastTime int64                // 客户端最后心跳时间/心跳检测
	closed   int32                // 客户端是否关闭连接
	channel  IChannel             // 渠道分组
	storage  cache.IBind          // 缓存服务
	event    IEvent               // 回调方法
	outChan  chan *ClientResponse // 发送通道
}

type ClientOption struct {
	Uid     int         // 用户识别ID
	Channel IChannel    // 渠道信息
	Storage cache.IBind // 自定义缓存组件, 用于绑定用户与客户端的关系
	Buffer  int         // 缓冲区大小根据业务, 自行调整
}

// ClientResponse 发送到客户端的消息结构体
type ClientResponse struct {
	IsAck   bool   `json:"-"`                 // 是否需要 ack 回调
	Sid     string `json:"sid,omitempty"`     // ACK ID
	Event   string `json:"event"`             // 事件名
	Content any    `json:"content,omitempty"` // 事件内容
	Retry   int    `json:"-"`                 // 重试次数(0 默认不重试)
}

// NewClient 初始化客户端
func NewClient(conn ISocket, option *ClientOption, event IEvent) error {
	if option.Buffer <= 0 {
		option.Buffer = 10
	}

	if event == nil {
		panic("event is nil")
	}

	client := &Client{
		conn:     conn,
		cid:      IdGen(),
		uid:      option.Uid,
		lastTime: time.Now().Unix(),
		channel:  option.Channel,
		storage:  option.Storage,
		event:    event,
		outChan:  make(chan *ClientResponse, option.Buffer),
	}

	// 设置客户端关闭回调事件
	conn.SetCloseHandler(client.hookClose)

	if option.Storage != nil {
		ctx := gctx.New()
		err := client.storage.Bind(ctx, client.channel.Name(), client.cid, client.uid)
		if err != nil {
			g.Log().Errorf(ctx, "bind client err:", err)
			return err
		}
	}

	// 注册客户端
	client.channel.addClient(client)

	// 触发自定义Open事件
	client.event.Open(client)

	// 注册心跳
	health.insert(client)

	return client.init()
}

func (c *Client) init() error {

	// 推送心跳机制配置
	_ = c.Write(&ClientResponse{
		Event: "heartbeat",
		Content: g.Map{
			"ping_interval": heartbeatInterval,
			"ping_timeout":  heartbeatTimeout,
		},
	})

	_ = grpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		c.loopWrite()
	}, nil)

	_ = grpool.AddWithRecover(gctx.New(), func(ctx context.Context) {
		c.loopAccept()
	}, nil)

	return nil
}

func (c *Client) Write(data *ClientResponse) error {
	defer func() {
		if err := recover(); err != nil {
			g.Log().Errorf(gctx.New(), "[%s-%d-%d] chan write err: %v", c.channel.Name(), c.cid, c.uid, err)
		}
	}()

	if c.Closed() {
		return gerror.New("connection has been closed")
	}

	if data.IsAck {
		data.Sid = fmt.Sprintf("%d_%d", IdGen(), gtime.Timestamp())
	}

	c.outChan <- data

	return nil
}

// Close 服务端主动关闭连接
func (c *Client) Close(code int, text string) {
	defer func() {
		if err := c.conn.Close(); err != nil {
			g.Log().Errorf(gctx.New(), "connection closed failed: %s", err.Error())
		}
	}()

	if err := c.hookClose(code, text); err != nil {
		g.Log().Errorf(gctx.New(), "[%s-%d-%d] client close err: %s", c.channel.Name(), c.cid, c.uid, err.Error())
	}
}

// Closed 获取客户端状态
func (c *Client) Closed() bool {
	return atomic.LoadInt32(&c.closed) == 1
}

func (c *Client) hookClose(code int, text string) error {
	// 在多协程并发环境确保只关闭一次，避免不必要的错误
	if !atomic.CompareAndSwapInt32(&c.closed, 0, 1) {
		return nil
	}

	// 关闭消息发送通道 outChan
	close(c.outChan)

	// 调用客户端关闭回调事件
	c.event.Close(c, code, text)

	//解绑用户和客户端关系
	if c.storage != nil {
		ctx := gctx.New()
		err := c.storage.UnBind(ctx, c.channel.Name(), c.cid)
		if err != nil {
			g.Log().Errorf(ctx, "unbind client err:", err)
			return err
		}
	}

	// 断开心跳监测
	health.delete(c)

	// 从通道中删除客户端
	c.channel.delClient(c)

	return nil
}

func (c *Client) Channel() IChannel {
	return c.channel
}

func (c *Client) Cid() int64 {
	return c.cid
}

func (c *Client) Uid() int {
	return c.uid
}

// 持续监听客户端，获取客户端发送的消息
func (c *Client) loopAccept() {
	defer c.Close(1000, "loop accept closed")

	for {
		_, data, err := c.conn.Read()
		if err != nil {
			g.Log().Errorf(gctx.New(), fmt.Sprintf("loop accept err: %s", err.Error()))
			break
		}

		c.lastTime = time.Now().Unix()

		c.handleMessage(data)
	}
}

// 监听发送通道，向客户端发送消息
func (c *Client) loopWrite() {
	ctx := gctx.New()

	timer := time.NewTimer(15 * time.Second)
	defer timer.Stop()

	for {
		timer.Reset(15 * time.Second)

		select {
		case <-timer.C:
			g.Log().Debugf(ctx, "client empty message cid: %d, uid: %d", c.cid, c.uid)
		case data, ok := <-c.outChan:
			if !ok || c.Closed() {
				return
			}

			jsonObj, err := gjson.Marshal(data)
			if err != nil {
				g.Log().Errorf(ctx, "client json marshal err: %s", err.Error())
				break
			}

			if err := c.conn.Write(jsonObj); err != nil {
				g.Log().Errorf(ctx, "[%s-%d-%d] client write err: %v", c.channel.Name(), c.cid, c.uid, err)
				return
			}

			if data.IsAck && data.Retry > 0 {
				data.Retry--

				ackContent := &AckContent{
					cid:      c.cid,
					uid:      gconv.Int64(c.uid),
					channel:  c.channel.Name(),
					response: data,
				}
				ack.insert(data.Sid, ackContent)
			}
		}
	}
}

// 处理接收到客户端的消息
func (c *Client) handleMessage(data []byte) {
	event, err := c.isJsonEvent(data)
	if err != nil {
		g.Log().Errorf(gctx.New(), "invalid data: ", err)
		return
	}

	switch event {
	case consts.MsgEventPing:
		_ = c.Write(&ClientResponse{Event: consts.MsgEventPong})
	case consts.MsgEventPong:
	case consts.MsgEventAck:
		j, _ := gjson.LoadJson(data)

		if !j.Get("sid").IsEmpty() {
			ack.delete(j.Get("sid").String())
		}
	default:
		c.event.Message(c, data)
	}
}

// 判断消息体是否符合要求
func (c *Client) isJsonEvent(data []byte) (string, error) {
	if !gjson.Valid(data) {
		return "", gerror.New("invalid json")
	}

	j, _ := gjson.LoadContent(data)
	if j.Get("event").IsEmpty() {
		return "", gerror.New("invalid event")
	}

	return j.Get("event").String(), nil
}
