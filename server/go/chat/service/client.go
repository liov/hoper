package service

import (
	"context"
	"flag"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hopeio/gox/encoding/json"
	"github.com/hopeio/protobuf/time/timestamp"
	"github.com/liov/hoper/server/go/chat/global"
	"github.com/liov/hoper/server/go/protobuf/chat"
)

var addr = flag.String("addr", "localhost:12345", "http service address")

// Client 表示一个 WebSocket 客户端连接
type Client struct {
	ID     uint64
	UID    uint64
	Conn   *websocket.Conn
	Send   chan []byte
	Hub    *Hub
	metaMu sync.RWMutex
}

// Write 将消息写入发送通道（非阻塞，满则跳过）
func (c *Client) Write(message []byte) bool {
	select {
	case c.Send <- message:
		return true
	default:
		return false
	}
}

// ReadPump 从连接读取消息并交给 Hub 处理
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister(c)
		c.Conn.Close()
	}()
	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// 可在此打日志
			}
			break
		}
		c.Hub.OnMessage(c, message)
	}
}

// WritePump 将 Send 通道中的消息写入连接
func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()
	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *Client) read(ctx context.Context) error {
	defer c.Hub.Unregister(c)
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				}
				break
			}
			var message chat.Message
			json.Unmarshal(msg, &message)
			message.CreatedAt = timestamp.New(time.Now())
			message.Uid = c.ID
			jsonMessage, _ := json.Marshal(&message)
			global.Dao.Redis.Do(ctx, "RPUSH", "Chat", jsonMessage)
		}
	}

}
