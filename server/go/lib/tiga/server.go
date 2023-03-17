package tiga

import (
	"context"
	"fmt"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"github.com/liov/hoper/server/go/lib/context/http_context"
	"github.com/liov/hoper/server/go/lib/initialize"
	"github.com/quic-go/quic-go/http3"
	"github.com/rs/cors"
	"go.opencensus.io/zpages"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/server/go/lib/utils/log"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"github.com/liov/hoper/server/go/lib/utils/net/http/grpc/gateway"
	"github.com/liov/hoper/server/go/lib/utils/strings"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

type CustomContext func(c context.Context, r *http.Request) context.Context
type ConvertContext func(r *http.Request) *http_context.Context

func (s *Server) Serve() {
	grpcServer := s.grpcHandler(s.Config)
	httpHandler := s.httpHandler(s.Config)

	// cors
	corsServer := cors.AllowAll()
	// grpc-web
	var wrappedGrpc *grpcweb.WrappedGrpcServer
	if s.Config.GrpcWeb {
		wrappedGrpc = NewGrpcWebServer(grpcServer)
	}

	openTracing := s.Config.OpenTracing
	//systemTracing := serviceConfig.SystemTracing
	if openTracing {
		grpc.EnableTracing = true
		/*opentracing.SetGlobalTracer(
		// tracing impl specific:
		basictracer.New(dapperish.NewTrivialRecorder(initialize.InitConfig.Module)),
		)*/
		//trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
		//trace.RegisterExporter(&exporter.PrintExporter{})
		zpages.Handle(http.DefaultServeMux, "/debug")
	}
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Default.Errorw(fmt.Sprintf("panic: %v", r), zap.String(log.Stack, stringsi.ToString(debug.Stack())))
				w.Header().Set(httpi.HeaderContentType, httpi.ContentJSONHeaderValue)
				w.Write(httpi.ResponseSysErr)
			}
		}()

		// 跨域
		if r.Method == http.MethodOptions && r.Header.Get(httpi.HeaderAccessControlRequestMethod) != "" {
			corsServer.HandlerFunc(w, r)
			return
		}

		ctx, span := http_context.ContextFromRequest(r, openTracing)
		if span != nil {
			defer span.End()
		}
		r = r.WithContext(ctx.ContextWrapper())

		contentType := r.Header.Get(httpi.HeaderContentType)
		if strings.HasPrefix(contentType, httpi.ContentGRPCHeaderValue) {
			if strings.HasPrefix(contentType[len(httpi.ContentGRPCHeaderValue):], "-web") && wrappedGrpc != nil {
				wrappedGrpc.ServeHTTP(w, r)
			} else if r.ProtoMajor == 2 && grpcServer != nil {
				grpcServer.ServeHTTP(w, r) // gRPC Server
			}
		} else {
			httpHandler(w, r)
		}
	})
	h2Handler := h2c.NewHandler(handle, new(http2.Server))
	server := &http.Server{
		Addr:         s.Config.Port,
		Handler:      h2Handler,
		ReadTimeout:  s.Config.ReadTimeout,
		WriteTimeout: s.Config.WriteTimeout,
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
		<-stop
		log.Debug("关闭服务")
		cs()
	}()

	if s.Config.Http3 != nil {
		go http3.ListenAndServe(s.Config.Http3.Address, s.Config.Http3.CertFile, s.Config.Http3.KeyFile, handle)
	}

	fmt.Println("listening: " + s.Config.Domain + s.Config.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Server struct {
	Config         *ServerConfig
	GRPCOptions    []grpc.ServerOption
	GRPCHandle     func(*grpc.Server)
	GatewayRegistr gateway.GatewayHandle
	GinHandle      func(engine *gin.Engine)
	GraphqlResolve graphql.ExecutableSchema
}

var stop = make(chan os.Signal, 1)

func (s *Server) Start() {
	if s.Config == nil {
		s.Config = defaultServerConfig()
	}
	signal.Notify(stop,
		// kill -SIGINT XXXX 或 Ctrl+c
		syscall.SIGINT, // register that too, it should be ok
		// os.Kill等同于syscall.Kill
		syscall.SIGKILL, // register that too, it should be ok
		// kill -SIGTERM XXXX
		syscall.SIGTERM,
	)
	s.Serve()
}

func NewServer(config *ServerConfig, ginhandle func(*gin.Engine), grpchandle func(*grpc.Server), grpcoptions []grpc.ServerOption, gatewayregist gateway.GatewayHandle, graphqlresolve graphql.ExecutableSchema) *Server {
	return &Server{
		Config:         config,
		GinHandle:      ginhandle,
		GRPCOptions:    grpcoptions,
		GRPCHandle:     grpchandle,
		GatewayRegistr: gatewayregist,
		GraphqlResolve: graphqlresolve,
	}
}

func Start(config *ServerConfig, ginhandle func(*gin.Engine), grpchandle func(*grpc.Server), grpcoptions []grpc.ServerOption, gatewayregist gateway.GatewayHandle, graphqlresolve graphql.ExecutableSchema) {
	NewServer(config, ginhandle, grpchandle, grpcoptions, gatewayregist, graphqlresolve).Start()
}
