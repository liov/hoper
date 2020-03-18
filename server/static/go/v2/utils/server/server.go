package server

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"runtime/debug"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/initialize"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var ch = make(chan os.Signal, 1)

func SignalChan() chan os.Signal {
	return ch
}

func (s *Server) Serve() {
	httpServer := s.Http()
	grpcServer := s.Grpc()
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.CallTwo.Error(" panic: ", r)
				debug.PrintStack()
				w.Write(errorcode.SysErr)
			}
		}()
		if r.ProtoMajor != 2 {
			httpServer.ServeHTTP(w, r)
			return
		}
		if strings.Contains(
			r.Header.Get("Content-Type"), "application/grpc",
		) {
			grpcServer.ServeHTTP(w, r) // gRPC Server
			return
		}

		httpServer.ServeHTTP(w, r)
		return
	})
	h2Handler := h2c.NewHandler(handle, &http2.Server{})
	//反射从配置中取port
	serverConfig := initialize.ServerConfig{Port: ":8080"}
	serverConfigType := reflect.TypeOf(&serverConfig).Elem()
	value := reflect.ValueOf(s.Conf).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).Type() == serverConfigType {
			serverConfig = value.Field(i).Interface().(initialize.ServerConfig)
		}
	}

	server := &http.Server{Addr: serverConfig.Port, Handler: h2Handler}
	go func() {
		log.Infof("listening%v", server.Addr)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	<-ch
	grpcServer.Stop()
	if err := server.Close(); err != nil {
		log.Error(err)
	}
}

type BasicConfig struct {
	initialize.BasicConfig
	Port string
}

type Config interface {
	initialize.Config
	GetBasicConfig() *initialize.BasicConfig
}

type Server struct {
	Conf        Config
	Dao         initialize.Dao
	GRPCRegistr func(*grpc.Server)
	HTTPRegistr func(context.Context, *runtime.ServeMux)
}

//
func (s *Server) Start() {
	defer v2.Start(s.Conf, s.Dao)()
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
			s.Serve()
		}
	}
}
