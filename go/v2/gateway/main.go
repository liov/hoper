package main

import (
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/kataras/iris"
	"github.com/kataras/iris/core/host"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	"github.com/liov/hoper/go/v2/gateway/internal/router"
	"github.com/liov/hoper/go/v2/initialize"
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
	flag.Parse()
	defer log.Sync()
	initialize.Start(config.Conf,nil)
	app := router.App()
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
			}), iris.WithConfiguration(iris.YAML("./config/iris.yml"))); err != nil && err != http.ErrServerClosed {
				log.Error(err)
			}
		}

	}
}
