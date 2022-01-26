package tiga

import (
	"context"
	"fmt"
	contexti "github.com/actliboy/hoper/server/go/lib/tiga/context"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/server"
	"go.opencensus.io/zpages"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/graphql"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/grpc/gateway"
	"github.com/actliboy/hoper/server/go/lib/utils/strings"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type CustomContext func(c context.Context, r *http.Request) context.Context
type ConvertContext func(r *http.Request) *contexti.Ctx

func (s *Server) Serve() {
	//反射从配置中取port
	serviceConfig := server.GetServiceConfig()
	grpcServer := s.grpcHandler(serviceConfig)
	httpHandler := s.httpHandler(serviceConfig)
	openTracing := serviceConfig.OpenTracing
	//systemTracing := serviceConfig.SystemTracing
	if openTracing {
		grpc.EnableTracing = true
		/*opentracing.SetGlobalTracer(
		// tracing impl specific:
		basictracer.New(dapperish.NewTrivialRecorder(initialize.InitConfig.Module)),
		)*/
		//trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		//trace.RegisterExporter(&exporter.PrintExporter{})
		zpages.Handle(http.DefaultServeMux, "/api/debug")
	}
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Default.Errorw(fmt.Sprintf("panic: %v", r), zap.String(log.Stack, stringsi.ToString(debug.Stack())))
				w.Header().Set(httpi.HeaderContentType, httpi.ContentJSONHeaderValue)
				w.Write(httpi.ResponseSysErr)
			}
		}()

		ctx, span := contexti.CtxFromRequest(r, openTracing)
		if span != nil {
			defer span.End()
		}
		r = r.WithContext(ctx.ContextWrapper())
		if r.ProtoMajor == 2 && grpcServer != nil && strings.Contains(r.Header.Get(httpi.HeaderContentType), httpi.ContentGRPCHeaderValue) {
			grpcServer.ServeHTTP(w, r) // gRPC Server
		} else {
			httpHandler(w, r)
		}
	})
	h2Handler := h2c.NewHandler(handle, new(http2.Server))
	server := &http.Server{
		Addr:         serviceConfig.Port,
		Handler:      h2Handler,
		ReadTimeout:  serviceConfig.ReadTimeout,
		WriteTimeout: serviceConfig.WriteTimeout,
	}
	// 服务注册
	initialize.InitConfig.Register()
	//服务关闭
	cs := func() {
		if grpcServer != nil {
			grpcServer.Stop()
		}
		if err := server.Close(); err != nil {
			log.Error(err)
		}
	}
	go func() {
		<-signals
		log.Debug("关闭服务")
		cs()
		signals <- syscall.SIGINT
	}()

	go func() {
		<-stop
		log.Debug("重启服务")
		cs()
	}()
	fmt.Println("listening: " + serviceConfig.Domain + serviceConfig.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Server struct {
	GRPCOptions    []grpc.ServerOption
	GRPCHandle     func(*grpc.Server)
	GatewayRegistr gateway.GatewayHandle
	GinHandle      func(engine *gin.Engine)
	GraphqlResolve graphql.ExecutableSchema
}

var signals = make(chan os.Signal, 1)
var stop = make(chan struct{}, 1)

func (s *Server) Start() {
	if initialize.InitConfig.EnvConfig == nil {
		log.Fatal(`初始化配置失败:
	main 函数的第一行应为
	defer initialize.Start(config.Conf, dao.Dao)()
`)
	}
	signal.Notify(signals,
		// kill -SIGINT XXXX 或 Ctrl+c
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	// 控制服务重启
Loop:
	for {
		select {
		case <-signals:
			break Loop
		default:
			s.Serve()
		}
	}
}

func ReStart() {
	stop <- struct{}{}
}
