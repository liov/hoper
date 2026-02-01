package service

import (
	"context"
	"sync"

	"google.golang.org/protobuf/proto"
)

var Manager *Hub

// Hub 管理所有 WebSocket 客户端连接
type Hub struct {
	clients   map[uint64]*Client
	mu        sync.RWMutex
	broadcast chan []byte
	peer      chan peerMsg
	serverID  uint16
}

type peerMsg struct {
	targetID uint64
	payload  []byte
}

// NewHub 创建 Hub，serverID 用于多实例时标识本节点
func NewHub(serverID uint16) *Hub {
	if serverID == 0 {
		panic("invalid server id")
	}
	h := &Hub{
		clients:   make(map[uint64]*Client),
		broadcast: make(chan []byte, 256),
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
		case payload := <-h.broadcast:
			h.broadcastToLocal(payload)
		case p := <-h.peer:
			h.sendToLocal(p.targetID, p.payload)
		}
	}
}

// BroadcastToLocal 仅向本机所有连接广播
func (h *Hub) broadcastToLocal(payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, c := range h.clients {
		c.Write(payload)
	}
}

// SendToLocal 仅向本机指定 clientID 发送
func (h *Hub) sendToLocal(clientID uint64, payload []byte) {
	h.mu.RLock()
	c, ok := h.clients[clientID]
	h.mu.RUnlock()
	if ok {
		c.Write(payload)
	}
}

// Broadcast 向所有客户端广播（若使用 Redis 则多实例都会收到）
func (h *Hub) Broadcast(payload []byte) {
	select {
	case h.broadcast <- payload:
	default:
		// 通道满可考虑丢弃或记录
	}
}

// SendToClient 向指定客户端发送（若使用 Redis 则仅对应实例会投递）
func (h *Hub) SendToClient(clientID uint64, payload []byte) {
	select {
	case h.peer <- peerMsg{targetID: clientID, payload: payload}:
	default:
	}
}

// OnMessage 由 Client.readPump 调用
func (h *Hub) OnMessage(c *Client, message []byte) {
	proto.Unmarshal(message)
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
