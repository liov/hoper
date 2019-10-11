package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	"github.com/kataras/iris"
	"github.com/liov/hoper/go/v2/gateway/router"
)

func main() {
	irisRouter := router.IrisRouter()
Loop:
	for {
		signal.Notify(router.Ch,
			// kill -SIGINT XXXX 或 Ctrl+c
			os.Interrupt,
			syscall.SIGINT, // register that too, it should be ok
			// os.Kill等同于syscall.Kill
			os.Kill,
			syscall.SIGKILL, // register that too, it should be ok
			// kill -SIGTERM XXXX
			syscall.SIGTERM,
		)
		select {
		case <-router.Ch:
			break Loop
		default:
			// listen and serve on https://0.0.0.0:8000.
			//if err := irisRouter.Run(iris.TLS(initialize.Config.Server.HttpPort, "../../config/tls/cert.pem", "../../config/tls/cert.key"),
			if err := irisRouter.Run(iris.Addr(config.Conf.Server.HttpPort),
				iris.WithConfiguration(iris.YAML("../../config/iris.yml"))); err != nil && err != http.ErrServerClosed {
				log.Error(err)
			}
		}

	}
	/*	opts := groupcache.HTTPPoolOptions{BasePath: hcache.BasePath}
		peers := groupcache.NewHTTPPoolOpts("", &opts)
		peers.Set("http://localhost:8333", "http://localhost:8222")

		val, err := hcache.GetFromPeer("helloworld", "wjs1", peers)*/
}
