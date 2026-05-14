package api

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	pb "github.com/liov/hoper/server/go/protobuf/remotebrowse"
	"github.com/liov/hoper/server/go/file/remotebrowse"
)

var rbWsUp = websocket.Upgrader{ReadBufferSize: 4096, WriteBufferSize: 4096, CheckOrigin: func(r *http.Request) bool { return true }}

func RegisterRemoteBrowse(app *gin.Engine) {
	pb.RegisterRemoteBrowseServiceHandlerServer(app, remotebrowse.GetService())
	app.GET("/rb/signal", remoteBrowseSignal)
}

func remoteBrowseSignal(c *gin.Context) {
	if upstream := strings.TrimSpace(os.Getenv("RB_SIGNAL_UPSTREAM")); upstream != "" {
		proxyRBSignal(c, upstream)
		return
	}
	hub := remotebrowse.EnsureHub()
	conn, err := rbWsUp.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	hub.HandleWS(conn)
}
