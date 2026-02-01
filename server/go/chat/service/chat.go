package service

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hopeio/gox/context/httpctx"
	"github.com/hopeio/gox/errors"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/liov/hoper/server/go/chat/global"
	"github.com/liov/hoper/server/go/protobuf/user"
)

const errResp = "未登录"

var manager = NewHub("1")

func Register(w http.ResponseWriter, r *http.Request) {
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
	ctxi, _ := httpctx.FromContext(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpx.CommonAnyResp{
			Code: errors.ErrCode(user.UserErrNoLogin),
			Msg:  errResp,
		}).ServeHTTP(w, r)
		return
	}
	client := &Client{ID: global.SF.Generate(), Conn: conn}

	manager.Register(client)

}
