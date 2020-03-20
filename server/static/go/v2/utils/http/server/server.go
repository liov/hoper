package server

import (
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/liov/hoper/go/v2/initialize"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/http/gateway"
	"github.com/liov/hoper/go/v2/utils/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

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
		if grpcServer != nil {
			if strings.Contains(
				r.Header.Get("Content-Type"), "application/grpc") {
				grpcServer.ServeHTTP(w, r) // gRPC Server
				return
			}
		}

		httpServer.ServeHTTP(w, r)
		return
	})
	h2Handler := h2c.NewHandler(handle, &http2.Server{})
	//反射从配置中取port

	server := &http.Server{Addr: getPort(s.Conf), Handler: h2Handler}
	cs := func() {
		if grpcServer != nil {
			grpcServer.Stop()
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

func getPort(v interface{}) string {
	value := reflect.ValueOf(v).Elem()
	for i := 0; i < value.NumField(); i++ {
		if conf, ok := value.Field(i).Interface().(initialize.ServerConfig); ok {
			return conf.Port

		}
	}
	return ":8080"
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
	HTTPRegistr gateway.GatewayHandle
}

var close = make(chan os.Signal, 1)
var stop = make(chan struct{}, 1)

func (s *Server) Start() {
	defer v2.Start(s.Conf, s.Dao)()
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
