package controller

import (
	"github.com/kataras/iris"
	"github.com/liov/hoper/go/v2/utils/log"
)

var info = struct {
	path     string
	version  uint8
	method   string
	describe string
	auth     string
	params   interface{}
}{}

type infoParam func()

func path(p string) infoParam {
	return func() {
		info.path = p
	}
}

func method(m string) infoParam {
	return func() {
		info.method = m
	}
}

func describe(d string) infoParam {
	return func() {
		info.describe = d
	}
}

func auth(a string) infoParam {
	return func() {
		info.auth = a
	}
}

func version(v uint8) infoParam {
	return func() {
		info.version = v
	}
}

func apiInfo(path,method,describe,auth,version,handle infoParam) {
	path()
	method()
	describe()
	auth()
	version()
	handle()
	log.Default.Infof("api: %s path:%s",info.describe,info.path)
}

func handle(app *iris.Application,h ...iris.Handler) infoParam {
	return func() {
		app.Handle(info.method,info.path,h...)
	}
}