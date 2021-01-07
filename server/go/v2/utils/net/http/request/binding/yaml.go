// Copyright 2018 Gin Core Team.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package binding

import (
	"bytes"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"gopkg.in/yaml.v2"
)

type yamlBinding struct{}

func (yamlBinding) Name() string {
	return "yaml"
}

func (yamlBinding) Bind(req *http.Request, obj interface{}) error {
	return decodeYAML(req.Body, obj)
}

func (yamlBinding) GinBind(ctx *gin.Context, obj interface{}) error {
	return decodeYAML(ctx.Request.Body, obj)
}

func (y yamlBinding) FasthttpBind(req *fasthttp.Request, obj interface{}) error {
	return y.BindBody(req.Body(), obj)
}

func (y yamlBinding) FiberBind(ctx *fiber.Ctx, obj interface{}) error {
	return y.BindBody(ctx.Body(), obj)
}

func (yamlBinding) BindBody(body []byte, obj interface{}) error {
	return decodeYAML(bytes.NewReader(body), obj)
}

func decodeYAML(r io.Reader, obj interface{}) error {
	decoder := yaml.NewDecoder(r)
	if err := decoder.Decode(obj); err != nil {
		return err
	}
	return validate(obj)
}
