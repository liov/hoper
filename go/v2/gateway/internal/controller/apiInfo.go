package controller

import (
	"github.com/kataras/iris"
	"github.com/liov/hoper/go/v2/utils/log"
)

type Controller struct {
	apiInfo
	App iris.Party
}

type apiInfo struct {
	path     string
	version  string
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

func version(v string) apiParam {
	return func(c *Controller) {
		c.version = v
	}
}

func handle(h ...iris.Handler) apiParam {
	return func(c *Controller) {
		c.App.Handle(c.apiInfo.method, c.path+"/v"+c.version, h...)
		log.Default.Infof("api: %s \t path:%s", c.apiInfo.describe, c.App.GetRelPath()+c.path +"/v"+c.version)
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