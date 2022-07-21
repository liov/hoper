// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"google.golang.org/protobuf/proto"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"

	"github.com/valyala/fasthttp"
)

type protobufBinding struct{}

func (protobufBinding) Name() string {
	return "protobuf"
}

func (b protobufBinding) Bind(req *http.Request, obj interface{}) error {
	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (b protobufBinding) GinBind(ctx *gin.Context, obj interface{}) error {
	buf, err := ioutil.ReadAll(ctx.Request.Body)
	if err != nil {
		return err
	}
	return b.BindBody(buf, obj)
}

func (b protobufBinding) FasthttpBind(req *fasthttp.Request, obj interface{}) error {
	return b.BindBody(req.Body(), obj)
}

func (b protobufBinding) FiberBind(ctx *fiber.Ctx, obj interface{}) error {
	return b.BindBody(ctx.Body(), obj)
}

func (protobufBinding) BindBody(body []byte, obj interface{}) error {
	if err := proto.Unmarshal(body, obj.(proto.Message)); err != nil {
		return err
	}
	// Here it's same to return validate(obj), but util now we can't add
	// `binding:""` to the struct which automatically generate by gen-proto
	return nil
	// return validate(obj)
}
