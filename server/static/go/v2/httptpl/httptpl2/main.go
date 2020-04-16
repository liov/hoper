package main

import (
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/httptpl/httptpl2/internal/config"
	_ "github.com/liov/hoper/go/v2/httptpl/httptpl2/internal/service"
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
	defer v2.Start(config.Conf, nil)()

	s := server.Server{
		IrisHandle: func(application *iris.Application) {
			api.RegisterAllService(application,
				initialize.Env == initialize.DEVELOPMENT,
				v2.BasicConfig.Module)
		},
	}
	s.Start()
}
