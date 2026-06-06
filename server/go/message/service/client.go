package service

import (
	"flag"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/hopeio/gox/log"
)

var addr = flag.String("addr", "localhost:12345", "http service address")

// Client 表示一个 WebSocket 客户端连接
type Client struct {
	ID      uint64
	UserID  uint64
	Device  string
	Channel string
	Conn    *websocket.Conn
	Send    chan msg
	Hub     *Hub
	metaMu  sync.RWMutex
}

// Write 将消息写入发送通道
func (c *Client) Write(msg msg) bool {
	select {
	case c.Send <- msg:
		return true
	default:
		return false
	}
}

// ReadPump 从连接读取消息并交给 Hub 处理
func (c *Client) ReadPump() {
	defer c.Hub.Unregister(c)
	for {
		typ, payload, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error()
			}
			break
		}
		c.Hub.OnMessage(c, typ, payload)
	}
}

// WritePump 将 Send 通道中的消息写入连接
func (c *Client) WritePump() {
	defer c.Hub.Unregister(c)
	for message := range c.Send {
		if message.PreparedMessage != nil {
			if err := c.Conn.WritePreparedMessage(message.PreparedMessage); err != nil {
				return
			}
		} else {
			if err := c.Conn.WriteMessage(message.typ, message.payload); err != nil {
				return
			}
		}
	}
}
