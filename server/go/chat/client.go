package chat

import (
	"flag"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/hopeio/pandora/utils/log"
)

var addr = flag.String("addr", "localhost:12345", "http service address")

func ClientStart() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}
	var dialer *websocket.Dialer

	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		log.Error(err)
		return
	}

	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Error("read:", err)
			return
		}

		log.Info("received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 2)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}
