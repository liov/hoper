package controller

import (
	"fmt"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/pio"
	"github.com/liov/hoper/go/v2/utils/strings2"
)

type Controller struct {
	*Handler
	Middle []iris.Handler
}
//时间换空间
type Handler struct {
	*ApiInfo
	App *iris.Application
}

type ApiInfo struct {
	path     string
	version  int
	method   string
	describe string
	auth     string
}

type apiParam func(c *Controller)

func path(p string) apiParam {
	return func(c *Controller) {
		c.path = p
	}
}

func method(m string) apiParam {
	return func(c *Controller) {
		c.method = m
	}
}

func describe(d string) apiParam {
	return func(c *Controller) {
		c.describe = d
	}
}

func auth(a string) apiParam {
	return func(c *Controller) {
		c.auth = a
	}
}

func version(v int) apiParam {
	return func(c *Controller) {
		c.version = v
	}
}

func handle(h ...iris.Handler) apiParam {
	return func(c *Controller) {
		path :="/api/v"+strconv.Itoa(c.version) + c.path
		handles := append(c.Middle,h...)
		c.App.Handle(c.ApiInfo.method, path, handles...)
		fmt.Printf(" %s\t %s %s\t %s\n",
			pio.Purple("API:"),
			pio.Yellow(strings2.FormatLen(c.ApiInfo.method,6)),
			pio.Blue(path,), pio.Gray(c.ApiInfo.describe))
	}
}

func (c *Controller) api(path, method, describe, auth, version, handle apiParam) {
	path(c)
	method(c)
	describe(c)
	auth(c)
	version(c)
	handle(c)
}