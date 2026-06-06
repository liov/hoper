package service

import (
	"context"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hopeio/gox/log"
	"github.com/hopeio/protobuf/time/timestamp"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/message"
	"google.golang.org/protobuf/proto"
)

var Manager *Hub

// Hub 管理所有 WebSocket 客户端连接
type Hub struct {
	clients   map[uint64]*Client
	mu        sync.RWMutex
	broadcast chan msg
	peer      chan peerMsg
	serverID  uint16
}

type msg struct {
	typ             int
	payload         []byte
	PreparedMessage *websocket.PreparedMessage
}
type peerMsg struct {
	targetID uint64
	msg
}

// NewHub 创建 Hub，serverID 用于多实例时标识本节点
func NewHub(serverID uint16) *Hub {
	if serverID == 0 {
		panic("invalid server id")
	}
	h := &Hub{
		clients:   make(map[uint64]*Client),
		broadcast: make(chan msg, 256),
		peer:      make(chan peerMsg, 256),
		serverID:  serverID,
	}
	return h
}

// Run 启动 Hub，处理广播/点对点
func (h *Hub) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-h.broadcast:
			h.broadcastToLocal(msg)
		case p := <-h.peer:
			h.sendToLocal(p.targetID, p.msg)
		}
	}
}

// BroadcastToLocal 仅向本机所有连接广播
func (h *Hub) broadcastToLocal(msg msg) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		c.Write(msg)
	}
}

// SendToLocal 仅向本机指定 clientID 发送
func (h *Hub) sendToLocal(clientID uint64, msg msg) {
	h.mu.RLock()
	c, ok := h.clients[clientID]
	h.mu.RUnlock()
	if ok {
		c.Write(msg)
	}
}

// Broadcast 向所有客户端广播（
func (h *Hub) Broadcast(typ int, payload []byte) error {
	preparedMessage, err := websocket.NewPreparedMessage(typ, payload)
	if err != nil {
		return err
	}
	select {
	case h.broadcast <- msg{
		PreparedMessage: preparedMessage,
	}:
	default:
		// 通道满可考虑丢弃或记录
	}
	return nil
}

// SendToClient 向指定客户端发送（若使用 Redis 则仅对应实例会投递）
func (h *Hub) SendToClient(clientID uint64, message msg) {
	select {
	case h.peer <- peerMsg{targetID: clientID, msg: message}:
	default:
	}
}

// OnMessage 由 Client.readPump 调用
func (h *Hub) OnMessage(c *Client, typ int, payload []byte) {
	var rmsg message.ClientMessage
	err := proto.Unmarshal(payload, &rmsg)
	if err != nil {
		log.Error(err)
		err = c.Conn.WriteMessage(websocket.TextMessage, []byte(err.Error()))
		if err != nil {
			log.Error(err)
			return
		}
		return
	}
	rmsg.ReadAt = timestamp.New(time.Now())
	_, err = global.MessageClient().Receive(context.TODO(), &rmsg)
	if err != nil {
		log.Error(err)
		return
	}
}

// Register 注册新客户端
func (h *Hub) Register(client *Client) {
	h.mu.Lock()
	h.clients[client.ID] = client
	defer h.mu.Unlock()
	go client.WritePump()
	go client.ReadPump()
}

// Unregister 注销客户端
func (h *Hub) Unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c.ID)
	close(c.Send)
	c.Conn.Close()
}

// ClientCount 当前本机连接数
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// GetClient 根据 ID 获取本机客户端（仅供本机逻辑使用）
func (h *Hub) GetClient(id uint64) *Client {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return h.clients[id]
}

// ServerID 返回本节点 ID
func (h *Hub) ServerID() uint16 {
	return h.serverID
}
