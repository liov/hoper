package controller

import (
	"github.com/kataras/iris"
	"github.com/liov/hoper/go/v2/utils/log"
)

type ApiInfo struct {
	path     string
	version  uint8
	method   string
	describe string
	auth     string
}

type infoParam func(info *ApiInfo)

func path(p string) infoParam {
	return func(info *ApiInfo) {
		info.path = p
	}
}

func method(m string) infoParam {
	return func(info *ApiInfo) {
		info.method = m
	}
}

func describe(d string) infoParam {
	return func(info *ApiInfo) {
		info.describe = d
	}
}

func auth(a string) infoParam {
	return func(info *ApiInfo) {
		info.auth = a
	}
}

func version(v uint8) infoParam {
	return func(info *ApiInfo) {
		info.version = v
	}
}

func (info *ApiInfo) apiInfo(path, method, describe, auth, version, handle infoParam) {
	path(info)
	method(info)
	describe(info)
	auth(info)
	version(info)
	handle(info)
	log.Default.Infof("api: %s path:%s", info.describe, info.path)
}

func handle(app *iris.Application, h ...iris.Handler) infoParam {
	return func(info *ApiInfo) {
		app.Handle(info.method, info.path, h...)
	}
}
