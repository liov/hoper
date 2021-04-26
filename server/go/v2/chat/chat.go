package chat

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/liov/hoper/go/v2/protobuf/user"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	"net/http"
)

const errRep = "上传失败"
const sep = "/"

func Chat(w http.ResponseWriter, r *http.Request) {
	conn, error := (&websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}).Upgrade(w, r, nil)
	if error != nil {
		http.NotFound(w, r)
		return
	}

	/*	if strings.Contains(c.Request().Header.Get("User-Agent"), "iPhone") {
			dviceName = "iPhone"
		} else if strings.Contains(c.Request().Header.Get("User-Agent"), "Android") {
			dviceName = "Android"
		} else {
			dviceName = "PC"
		}*/
	ctxi := contexti.CtxFromContext(r.Context())
	_, err := auth(ctxi, false)
	if err != nil {
		(&httpi.ResData{
			Code:    uint32(user.UserErrLogin),
			Message: errRep,
		}).Response(w)
		return
	}
	client := &Client{uuid: uuid.New().String(), conn: conn, send: make(chan []byte), ctx: ctxi}

	manager.register <- client

	go client.read()
	go client.write()
}
