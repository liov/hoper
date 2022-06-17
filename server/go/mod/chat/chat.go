package chat

import (
	"github.com/actliboy/hoper/server/go/lib/protobuf/errorcode"
	contexti "github.com/actliboy/hoper/server/go/lib/tiga/context"
	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/actliboy/hoper/server/go/mod/protobuf/user"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
)

const errRep = "未登录"

func init() {
	go manager.start()
}

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
			Code:    errorcode.ErrCode(user.UserErrLogin),
			Message: errRep,
		}).Response(w)
		return
	}

	client := &Client{uuid: uuid.New().String(), conn: conn, send: make(chan []byte), ctx: ctxi}

	manager.register <- client

	go client.read()
	go client.write()
}
