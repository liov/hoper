package controller

import (
	"github.com/kataras/iris"
	"github.com/liov/hoper/go/v2/utils/log"
)

type apiInfo struct {
	path     string
	version  uint8
	method   string
	describe string
	auth     string
}

type infoParam func(info *apiInfo)

func path(p string) infoParam {
	return func(info *apiInfo) {
		info.path = p
	}
}

func method(m string) infoParam {
	return func(info *apiInfo) {
		info.method = m
	}
}

func describe(d string) infoParam {
	return func(info *apiInfo) {
		info.describe = d
	}
}

func auth(a string) infoParam {
	return func(info *apiInfo) {
		info.auth = a
	}
}

func version(v uint8) infoParam {
	return func(info *apiInfo) {
		info.version = v
	}
}

func (info *apiInfo) api(path, method, describe, auth, version, handle infoParam) {
	path(info)
	method(info)
	describe(info)
	auth(info)
	version(info)
	handle(info)
	log.Default.Infof("api: %s path:%s", info.describe, info.path)
}

func handle(app *iris.Application, h ...iris.Handler) infoParam {
	return func(info *apiInfo) {
		app.Handle(info.method, info.path, h...)
	}
}
