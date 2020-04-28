package server

import (
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/graphql"
	"github.com/liov/hoper/go/v2/initialize"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Server) Serve() {
	httpServer := s.Http()

	if s.GRPCServer != nil {
		reflection.Register(s.GRPCServer)
		v2.BasicConfig.Register()
	}

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
		if s.GRPCServer != nil {
			if strings.Contains(
				r.Header.Get("Content-Type"), "application/grpc") {
				s.GRPCServer.ServeHTTP(w, r) // gRPC Server
				return
			}
		}

		httpServer.ServeHTTP(w, r)
		return
	})
	h2Handler := h2c.NewHandler(handle, &http2.Server{})
	//反射从配置中取port
	serviceConfig := v2.BasicConfig.GetServiceConfig()
	server := &http.Server{
		Addr:         serviceConfig.Port,
		Handler:      h2Handler,
		ReadTimeout:  serviceConfig.ReadTimeout,
		WriteTimeout: serviceConfig.WriteTimeout,
	}
	cs := func() {
		if s.GRPCServer != nil {
			s.GRPCServer.Stop()
		}
		if err := server.Close(); err != nil {
			log.Error(err)
		}
	}
	go func() {
		<-close
		log.Info("关闭服务")
		cs()
		close <- syscall.SIGINT
	}()

	go func() {
		<-stop
		log.Info("重启服务")
		cs()
	}()
	log.Infof("listening%v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Config interface {
	initialize.Config
}

type Server struct {
	GRPCServer     *grpc.Server
	GatewayRegistr gateway.GatewayHandle
	PickHandle     func(*pick.EasyRouter)
	GraphqlResolve graphql.ExecutableSchema
}

var close = make(chan os.Signal, 1)
var stop = make(chan struct{}, 1)

func (s *Server) Start() {
	if v2.BasicConfig == nil {
		log.Fatal(`初始化配置失败:
	main 函数的第一行应为
	defer v2.Start(config.Conf, dao.Dao)()
`)
	}
	signal.Notify(close,
		// kill -SIGINT XXXX 或 Ctrl+c
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
Loop:
	for {
		select {
		case <-close:
			break Loop
		default:
			s.Serve()
		}
	}
}

func ReStart() {
	stop <- struct{}{}
}
