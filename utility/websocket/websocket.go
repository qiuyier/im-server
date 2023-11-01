// Package websocket 处理 WebSocket 连接，提供了一个通用的 WebSocket 连接管理接口，并使用 Gorilla WebSocket 包作为底层实现
package websocket

import (
	"github.com/gorilla/websocket"
	"net/http"
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

// WsConn 基于 Gorilla WebSocket 实现 ISocket
type WsConn struct {
	conn *websocket.Conn
}

// 配置 WebSocket 连接的参数
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// NewWebSocket 创建一个 WebSocket 连接
func NewWebSocket(w http.ResponseWriter, r *http.Request) (*WsConn, error) {
	conn, err := upgrader.Upgrade(w, r, w.Header())
	if err != nil {
		return nil, err
	}

	return &WsConn{conn: conn}, nil
}

func (w *WsConn) Read() (int, []byte, error) {
	return w.conn.ReadMessage()
}

func (w *WsConn) Write(bytes []byte) error {
	return w.conn.WriteMessage(websocket.TextMessage, bytes)
}

func (w *WsConn) Close() error {
	return w.conn.Close()
}

func (w *WsConn) SetCloseHandler(h func(code int, text string) error) {
	w.conn.SetCloseHandler(h)
}
