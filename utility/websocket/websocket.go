// Package websocket 处理 WebSocket 连接，提供了一个通用的 WebSocket 连接管理接口，并使用 Gorilla WebSocket 包作为底层实现
package websocket

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

// ISocket 定义 WebSocket 连接的通用操作。它包括读数据、写数据、关闭连接以及设置连接关闭回调事件。
// 用于提供 WebSocket 连接的通用行为，以便可以有不同的底层实现
type ISocket interface {
	// Read 读数据
	Read() (int, []byte, error)

	// Write 写数据
	Write(bytes []byte) error

	// Close 连接关闭
	Close() error

	// SetCloseHandler 设置连接关闭回调事件
	SetCloseHandler(h func(code int, text string) error)
}

// GfConn GoFrame已封装好Gorilla，直接使用即可
type GfConn struct {
	conn *ghttp.WebSocket
}

func NewGfWebSocket(r *ghttp.Request) (*GfConn, error) {
	conn, err := r.WebSocket()
	if err != nil {
		return nil, err
	}

	return &GfConn{conn: conn}, nil
}

func (w *GfConn) Read() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *GfConn) Write(bytes []byte) error {
	return w.conn.WriteMessage(ghttp.WsMsgText, bytes)
}

func (w *GfConn) Close() error {
	return w.conn.Close()
}

func (w *GfConn) SetCloseHandler(h func(code int, text string) error) {
	w.conn.SetCloseHandler(h)
}
