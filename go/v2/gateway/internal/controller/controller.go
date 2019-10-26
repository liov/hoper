package controller

import (
	"strconv"

	"github.com/kataras/iris"
	"github.com/liov/hoper/go/v2/utils/log"
)

type Controller struct {
	apiInfo
	App *iris.Application
	Middle []iris.Handler
}

type apiInfo struct {
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
		c.App.Handle(c.apiInfo.method, path, handles...)
		log.NoCall.Infof("method:%s api:%s \t uri:%s", c.apiInfo.method, c.apiInfo.describe, path)
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