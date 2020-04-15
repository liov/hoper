package main

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/liov/hoper/go/v2/httptpl3/internal/config"
	"github.com/liov/hoper/go/v2/httptpl3/internal/service"
	_ "github.com/liov/hoper/go/v2/httptpl3/internal/service"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
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
		IrisHandle: func(app *iris.Application) {
			svc := &service.UserService{}
			mvc.New(app).Register(mvc.AutoBinding).Handle(svc)
			app.DI().Handle(iris.MethodPost, "/{id:int}", svc.Add)
		},
	}
	s.Start()
}
