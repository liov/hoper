package service

import (
	"github.com/gorilla/websocket"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/encoding/json"
	"github.com/liov/hoper/server/go/content/confdao"
	"github.com/liov/hoper/server/go/protobuf/user"
	"sync"
	"time"
)

type ClientManager struct {
	mutex   sync.Mutex
	clients map[uint64]*Client
}

type Client struct {
	id   uint64
	conn *websocket.Conn
	ctx  *httpctx.Context
}

type Message struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	SendUserID uint64    `json:"sendUserId,omitempty"`
	RecvUserID uint64    `json:"recvUserId,omitempty"`
	Remark     string    `json:"remark,omitempty"`
	Device     string    `json:"device"`
	Content    string    `json:"content,omitempty"`
}

var manager = ClientManager{
	clients: make(map[uint64]*Client),
}

func (manager *ClientManager) register(client *Client) {
	manager.mutex.Lock()
	id := client.ctx.AuthInfo.(*user.AuthInfo).Id
	manager.clients[id] = client
	jsonMessage, _ := json.Marshal(&Message{Remark: "/A new conn has connected."})
	client.conn.WriteMessage(websocket.PongMessage, jsonMessage)
	manager.mutex.Unlock()
}

func (manager *ClientManager) unregister(client *Client) {
	manager.mutex.Lock()
	id := client.ctx.AuthInfo.(*user.AuthInfo).Id
	if _, ok := manager.clients[id]; ok {
		delete(manager.clients, id)
		jsonMessage, _ := json.Marshal(&Message{Remark: "/A conn has disconnected."})
		client.conn.WriteMessage(websocket.CloseMessage, jsonMessage)
	}
	client.conn.Close()
	manager.mutex.Unlock()
}

func (manager *ClientManager) broadcast(message []byte) {
	for _, client := range manager.clients {
		client.conn.WriteMessage(websocket.TextMessage, message)
	}
}

func (c *Client) read() {
	defer manager.unregister(c)
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}
		var message Message
		json.Unmarshal(msg, &message)
		message.CreatedAt = time.Now()
		message.SendUserID = c.ctx.AuthInfo.(*user.AuthInfo).Id
		jsonMessage, _ := json.Marshal(&message)
		confdao.Dao.Redis.Do(c.ctx.BaseContext(), "RPUSH", "Chat", jsonMessage)
	}

}
