package service

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/hopeio/context/httpctx"
	"github.com/hopeio/gox/errors"
	"github.com/hopeio/gox/idgen/snowflake"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/liov/hoper/server/go/protobuf/user"
)

const errRep = "未登录"

var idgen = snowflake.NewSnowflake(snowflake.Settings{})

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
	ctxi, _ := httpctx.FromContext(r.Context())
	_, err = auth(ctxi, false)
	if err != nil {
		(&httpx.RespAnyData{
			Code: errors.ErrCode(user.UserErrNoLogin),
			Msg:  errRep,
		}).Response(w)
		return
	}
	id, _ := idgen.NextID()
	client := &Client{id: id, conn: conn, ctx: ctxi}

	manager.register(client)

	go client.read()
}
