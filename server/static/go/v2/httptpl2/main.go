package main

import (
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/httptpl2/internal/config"
	"github.com/liov/hoper/go/v2/httptpl2/internal/router"
	"github.com/liov/hoper/go/v2/initialize"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
)

func main() {
	/*	f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		trace.Start(f)
		defer trace.Stop()*/
	s := server.Server{
		Conf: config.Conf,
		IrisHandle: func(application *iris.Application) {
			api.RegisterAllService(application, router.Route(),
				initialize.Env == initialize.DEVELOPMENT,
				v2.BasicConfig.Module)
		},
	}
	s.Start()
}
