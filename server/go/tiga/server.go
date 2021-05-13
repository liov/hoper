package tiga

import (
	"context"
	"encoding/base64"
	"fmt"
	contexti "github.com/liov/hoper/v2/tiga/context"
	"net/http"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/graphql"
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/v2/tiga/initialize"
	"github.com/liov/hoper/v2/utils/log"
	httpi "github.com/liov/hoper/v2/utils/net/http"
	"github.com/liov/hoper/v2/utils/net/http/grpc/gateway"
	"github.com/liov/hoper/v2/utils/strings"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	gtrace "golang.org/x/net/trace"
	"google.golang.org/grpc"
)

type CustomContext func(c context.Context, r *http.Request) context.Context
type ConvertContext func(r *http.Request) *contexti.Ctx

func (s *Server) Serve() {
	//反射从配置中取port
	serviceConfig := initialize.InitConfig.GetServiceConfig()
	grpcServer := s.grpcHandler(serviceConfig)
	httpHandler := s.httpHandler(serviceConfig)
	openTracing := serviceConfig.OpenTracing
	systemTracing := serviceConfig.SystemTracing
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.Default.Errorw(fmt.Sprintf("panic: %v", r), zap.String(log.Stack, stringsi.ToString(debug.Stack())))
				w.Header().Set(httpi.HeaderContentType, httpi.ContentJSONHeaderValue)
				w.Write(httpi.ResponseSysErr)
			}
		}()

		// 请求TraceID，链路跟踪用
		ctx := r.Context()
		if systemTracing {
			// 系统trace只能追踪单个请求，且只记录时间及是否完成
			t := gtrace.New(initialize.InitConfig.Module, r.RequestURI)
			defer t.Finish()
			ctx = gtrace.NewContext(ctx, t)
		}
		if openTracing {
			var span *trace.Span
			// 直接从远程读取Trace信息，Trace是否为空交给propagation包判断
			traceString := r.Header.Get(httpi.GrpcTraceBin)
			var traceBin []byte
			if len(traceString)%4 == 0 {
				// Input was padded, or padding was not necessary.
				traceBin, _ = base64.StdEncoding.DecodeString(traceString)
			}
			traceBin, _ = base64.RawStdEncoding.DecodeString(traceString)
			if parent, ok := propagation.FromBinary(traceBin); ok {
				ctx, span = trace.StartSpanWithRemoteParent(ctx, r.RequestURI,
					parent, trace.WithSampler(trace.AlwaysSample()),
					trace.WithSpanKind(trace.SpanKindServer))
			} else {
				ctx, span = trace.StartSpan(ctx, r.RequestURI,
					trace.WithSampler(trace.AlwaysSample()),
					trace.WithSpanKind(trace.SpanKindServer))
			}
			defer span.End()
		}

		ctx = contexti.CtxWithRequest(ctx, r).ContextWrapper()

		r = r.WithContext(ctx)
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
