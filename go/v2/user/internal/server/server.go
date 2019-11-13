package server

import (
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/utils/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

var ch = make(chan os.Signal, 1)

func SignalChan() chan os.Signal {
	return ch
}

func Serve() {
	signal.Notify(ch,
		// kill -SIGINT XXXX 或 Ctrl+c
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor != 2 {
			Http().ServeHTTP(w, r)
			return
		}
		if strings.Contains(
			r.Header.Get("Content-Type"), "application/grpc",
		) {
			Grpc().ServeHTTP(w, r) // gRPC Server
			return
		}

		Http().ServeHTTP(w, r)
		return
	})
	h2Handler := h2c.NewHandler(handle, &http2.Server{})
	server := &http.Server{Addr: config.Conf.Server.Port, Handler: h2Handler}
	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ch
	Grpc().Stop()
	server.Close()
}
