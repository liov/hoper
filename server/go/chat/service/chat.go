package service

import (
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/utils/errors/errcode"
	httpi "github.com/hopeio/utils/net/http"
	"github.com/liov/hoper/server/go/protobuf/user"
	"net/http"
)

const errRep = "未登录"

func init() {
	go manager.start()
}

func Chat(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}).Upgrade(w, r, nil)
	if err != nil {
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
	ctxi := httpctx.FromContextValue(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpi.ResAnyData{
			Code:    errcode.ErrCode(user.UserErrNoLogin),
			Message: errRep,
		}).Response(w, http.StatusOK)
		return
	}

	client := &Client{uuid: uuid.New().String(), conn: conn, send: make(chan []byte), ctx: ctxi}

	manager.register <- client

	go client.read()
	go client.write()
}
