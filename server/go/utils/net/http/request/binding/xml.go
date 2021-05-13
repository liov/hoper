// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"encoding/xml"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type xmlBinding struct{}

func (xmlBinding) Name() string {
	return "xml"
}

func (xmlBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeXML(req.Body, obj)
}

func (xmlBinding) GinBind(ctx *gin.Context, obj interface{}) error {
	return decodeXML(ctx.Request.Body, obj)
}

func (x xmlBinding) FasthttpBind(req *fasthttp.Request, obj interface{}) error {
	return x.BindBody(req.Body(), obj)
}

func (x xmlBinding) FiberBind(ctx *fiber.Ctx, obj interface{}) error {
	return x.BindBody(ctx.Body(), obj)
}

func (xmlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeXML(bytes.NewReader(body), obj)
}
func decodeXML(r io.Reader, obj interface{}) error {
	decoder := xml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
