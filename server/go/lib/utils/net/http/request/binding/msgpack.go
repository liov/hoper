// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

// +build !nomsgpack

package binding

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/ugorji/go/codec"
	"github.com/valyala/fasthttp"
)

type msgpackBinding struct{}

func (msgpackBinding) Name() string {
	return "msgpack"
}

func (msgpackBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeMsgPack(req.Body, obj)
}

func (msgpackBinding) GinBind(ctx *gin.Context, obj interface{}) error {
	return decodeMsgPack(ctx.Request.Body, obj)
}

func (m msgpackBinding) FasthttpBind(req *fasthttp.Request, obj interface{}) error {
	return m.BindBody(req.Body(), obj)
}

func (m msgpackBinding) FiberBind(ctx *fiber.Ctx, obj interface{}) error {
	return m.BindBody(ctx.Body(), obj)
}

func (msgpackBinding) BindBody(body []byte, obj interface{}) error {
	return decodeMsgPack(bytes.NewReader(body), obj)
}

func decodeMsgPack(r io.Reader, obj interface{}) error {
	cdc := new(codec.MsgpackHandle)
	if err := codec.NewDecoder(r, cdc).Decode(&obj); err != nil {
		return err
	}
	return validate(obj)
}
