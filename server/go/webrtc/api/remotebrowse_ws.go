package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/liov/hoper/server/go/webrtc/service"
)

var wsUp = websocket.Upgrader{ReadBufferSize: 4096, WriteBufferSize: 4096, CheckOrigin: func(r *http.Request) bool { return true }}

func remoteBrowseSignal(c *gin.Context) {
	hub := service.EnsureSignalHub()
	conn, err := wsUp.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	hub.HandleWS(conn)
}
