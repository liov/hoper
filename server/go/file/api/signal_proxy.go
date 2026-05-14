package api

import (
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func proxyRBSignal(c *gin.Context, upstream string) {
	u, err := url.Parse(upstream)
	if err != nil {
		c.String(http.StatusBadGateway, "bad RB_SIGNAL_UPSTREAM")
		return
	}
	if u.Scheme == "http" {
		u.Scheme = "ws"
	}
	if u.Scheme == "https" {
		u.Scheme = "wss"
	}
	clientConn, err := rbWsUp.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer clientConn.Close()
	serverConn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return
	}
	defer serverConn.Close()
	errCh := make(chan error, 2)
	go func() { errCh <- copyWS(clientConn, serverConn) }()
	go func() { errCh <- copyWS(serverConn, clientConn) }()
	<-errCh
}

func copyWS(dst, src *websocket.Conn) error {
	for {
		mt, data, err := src.ReadMessage()
		if err != nil {
			return err
		}
		if err := dst.WriteMessage(mt, data); err != nil {
			return err
		}
	}
}
