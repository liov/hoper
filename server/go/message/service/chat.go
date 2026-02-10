package service

import (
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/hopeio/gox/errors"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/scaffold/errcode"
	"github.com/liov/hoper/server/go/protobuf/user"
)

const errResp = "未登录"

func Register(w http.ResponseWriter, r *http.Request) {
	conn, err := (&websocket.Upgrader{
		CheckOrigin:     func(r *http.Request) bool { return true },
		ReadBufferSize:  1024,
		WriteBufferSize: 1024}).Upgrade(w, r, nil)
	if err != nil {
		(&httpx.CommonAnyResp{
			Code: errors.ErrCode(errcode.InvalidArgument),
			Msg:  err.Error(),
		}).ServeHTTP(w, r)
		return
	}

	device := "PC"
	channel := "wx"
	if strings.Contains(r.Header.Get("User-Agent"), "iPhone") {
		device = "iPhone"
	} else if strings.Contains(r.Header.Get("User-Agent"), "Android") {
		device = "Android"
	} else {
		device = "PC"
	}
	if strings.Contains(r.Header.Get("User-Agent"), "MicroMessenger") {
		channel = "wx"
	}
	ctx, span := Tracer.Start(r.Context(), "Register")
	defer span.End()
	auth, err := auth(ctx, false)
	if err != nil {
		(&httpx.CommonAnyResp{
			Code: errors.ErrCode(user.UserErrNoLogin),
			Msg:  errResp,
		}).ServeHTTP(w, r)
		return
	}
	client := &Client{ID: auth.Id, Conn: conn, Device: device, Channel: channel}
	Manager.Register(client)
}
