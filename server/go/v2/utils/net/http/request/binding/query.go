// Copyright 2017 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

type queryBinding struct{}

func (queryBinding) Name() string {
	return "query"
}

func (queryBinding) Bind(req *http.Request, obj interface{}) error {
	values := req.URL.Query()
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

func (queryBinding) GinBind(ctx *gin.Context, obj interface{}) error {
	values := ctx.Request.URL.Query()
	for i:=range ctx.Params{
		values.Set(ctx.Params[i].Key,ctx.Params[i].Value)
	}
	if err := mapForm(obj, values); err != nil {
		return err
	}
	return validate(obj)
}

func (queryBinding) FasthttpBind(req *fasthttp.Request, obj interface{}) error {
	values:=req.URI().QueryArgs()
	if err := mapKV(obj, (*argsSource)(values)); err != nil {
		return err
	}
	return validate(obj)
}

func (queryBinding) FiberBind(ctx *fiber.Ctx, obj interface{}) error {
	values:=ctx.Request().URI().QueryArgs()
	if err := mapKV(obj, (*argsSource)(values)); err != nil {
		return err
	}
	return validate(obj)
}
