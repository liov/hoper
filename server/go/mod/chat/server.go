package chat

import (
	"github.com/gorilla/websocket"
	"github.com/liov/hoper/server/go/lib/context/http_context"
	"github.com/liov/hoper/server/go/lib/utils/encoding/json/iterator"
	"github.com/liov/hoper/server/go/mod/content/dao"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
	"time"
)

type ClientManager struct {
	clients    map[uint64]*Client
	broadcast  chan []byte //广播聊天
	register   chan *Client
	unregister chan *Client
}

type Client struct {
	uuid string
	conn *websocket.Conn
	send chan []byte
	ctx  *http_context.Context
}

type Message struct {
	ID         uint64    `gorm:"primary_key" json:"id"`
	CreatedAt  time.Time `json:"createdAt"`
	SendUserID uint64    `json:"sendUserId,omitempty"`
	RecvUserID uint64    `json:"recvUserId,omitempty"`
	Content    string    `json:"content,omitempty"`
	Remark     string    `json:"remark,omitempty"`
	Device     string    `json:"device"`
}

var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[uint64]*Client),
}

func (manager *ClientManager) start() {
	for {
		select {
		case client := <-manager.register:
			id := client.ctx.Props.AuthInfo.(*user.AuthInfo).Id
			manager.clients[id] = client
			jsonMessage, _ := iterator.Marshal(&Message{Remark: "/A new conn has connected."})
			manager.send(jsonMessage, client)
		case client := <-manager.unregister:
			id := client.ctx.Props.AuthInfo.(*user.AuthInfo).Id
			if _, ok := manager.clients[id]; ok {
				close(client.send)
				delete(manager.clients, id)
				jsonMessage, _ := iterator.Marshal(&Message{Remark: "/A conn has disconnected."})
				manager.send(jsonMessage, client)
			}
		case message := <-manager.broadcast:
			//这里貌似可以做单点发送
			for _, client := range manager.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					id := client.ctx.Props.AuthInfo.(*user.AuthInfo).Id
					delete(manager.clients, id)
				}
			}
		}
	}
}

func (manager *ClientManager) send(message []byte, ignore *Client) {
	for _, client := range manager.clients {
		if client != ignore {
			client.send <- message
		}
	}
}

func (c *Client) read() {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
			}
			break
		}
		var message Message
		iterator.Unmarshal(msg, &message)
		message.CreatedAt = time.Now()
		message.SendUserID = c.ctx.Props.AuthInfo.(*user.AuthInfo).Id
		jsonMessage, _ := iterator.Marshal(&message)
		dao.Dao.Redis.Do(c.ctx, "RPUSH", "Chat", jsonMessage)
		manager.broadcast <- jsonMessage
	}
}

func (c *Client) write() {
	defer func() {
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			c.conn.WriteMessage(websocket.TextMessage, message)
		}
	}
}
