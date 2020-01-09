package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/host"
	"github.com/liov/hoper/go/v2/httptpl/internal/config"
	"github.com/liov/hoper/go/v2/httptpl/internal/router"
	"github.com/liov/hoper/go/v2/initialize"
	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/log"
)

func main() {
	/*	f, err := os.Create("trace.out")
		if err != nil {
			panic(err)
		}
		defer f.Close()

		trace.Start(f)
		defer trace.Stop()*/
	defer initialize.Start(config.Conf, nil)()
	app := router.App()
	iris_build.WithConfiguration(app, initialize.ConfUrl)
	ch := make(chan os.Signal, 1)
Loop:
	for {
		signal.Notify(ch,
			// kill -SIGINT XXXX 或 Ctrl+c
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-ch:
			break Loop
		default:
			// listen and serve on https://0.0.0.0:8000.
			//if err := irisRouter.Run(iris.TLS(initialize.config.Server.HttpPort, "../../config/tls/cert.pem", "../../config/tls/cert.key"),
			if err := app.Run(iris.Addr(config.Conf.Server.Port, func(su *host.Supervisor) {
				su.Server.WriteTimeout = config.Conf.Server.WriteTimeout
				su.Server.ReadTimeout = config.Conf.Server.ReadTimeout
			})); err != nil && err != http.ErrServerClosed {
				log.Error(err)
			}
		}

	}
}
