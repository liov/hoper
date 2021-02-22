package tailmon

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/log"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	gin_build "github.com/liov/hoper/go/v2/utils/net/http/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	runtimei "github.com/liov/hoper/go/v2/utils/runtime"
	"github.com/liov/hoper/go/v2/utils/strings"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	gtrace "golang.org/x/net/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Server) httpHandler() http.HandlerFunc {
	// 默认使用gin
	ginServer := gin_build.Http(initialize.InitConfig.ConfUrl, s.GinHandle)

	if s.GraphqlResolve != nil {
		graphqlServer := handler.NewDefaultServer(s.GraphqlResolve)
		ginServer.Handle(http.MethodPost, "/api/graphql", func(ctx *gin.Context) {
			graphqlServer.ServeHTTP(ctx.Writer, ctx.Request)
		})
	}
	var gatewayServer http.Handler
	if s.GatewayRegistr != nil {
		gatewayServer = gateway.Gateway(s.GatewayRegistr)
		/*	ginServer.NoRoute(func(ctx *gin.Context) {
			gatewayServer.ServeHTTP(
				(*httpi.ResponseRecorder)(unsafe.Pointer(uintptr(*(*int64)(unsafe.Pointer(uintptr(unsafe.Pointer(ctx))+8))))),
				ctx.Request)
			ctx.Writer.WriteHeader(http.StatusOK)
		})*/
	}

	// http.Handle("/", ginServer)
	var excludes = []string{"/debug", "/api-doc", "/metrics"}
	return func(w http.ResponseWriter, r *http.Request) {
		// 暂时解决方法，三个路由
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			h.ServeHTTP(w, r)
		}
		if stringsi.HasPrefixes(r.RequestURI, excludes) {
			ginServer.ServeHTTP(w, r)
			return
		}

		var body []byte
		if r.Method != http.MethodGet {
			body, _ = ioutil.ReadAll(r.Body)
			r.Body = ioutil.NopCloser(bytes.NewReader(body))
		}
		recorder := httpi.NewRecorder(w.Header())

		ginServer.ServeHTTP(recorder, r)
		if recorder.Code == http.StatusNotFound && gatewayServer != nil {
			recorder.Reset()
			gatewayServer.ServeHTTP(recorder, r)
		}

		// 提取 recorder 中记录的状态码，写入到 ResponseWriter 中
		w.WriteHeader(recorder.Code)
		if recorder.Body != nil {
			// 将 recorder 记录的 Response Body 写入到 ResponseWriter 中，客户端收到响应报文体
			w.Write(recorder.Body.Bytes())
		}

		accessLog(s.ConvertContext(r), r.RequestURI,
			stringsi.ToString(body), stringsi.ToString(recorder.Body.Bytes()),
			recorder.Code)
	}
}

type CustomContext func(c context.Context, r *http.Request) context.Context
type ConvertContext func(r *http.Request) pick.Context

func (s *Server) Serve() {
	//反射从配置中取port
	serviceConfig := initialize.InitConfig.GetServiceConfig()
	var grpcServer *grpc.Server
	if s.GRPCHandle != nil {
		if serviceConfig.Prometheus {
			s.GRPCOptions = append([]grpc.ServerOption{
				grpc.ChainStreamInterceptor(grpc_prometheus.StreamServerInterceptor),
				grpc.ChainUnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
			}, s.GRPCOptions...)
		}
		grpcServer = grpc.NewServer(s.GRPCOptions...)
		if serviceConfig.Prometheus {
			grpc_prometheus.Register(grpcServer)
		}
		s.GRPCHandle(grpcServer)
		reflection.Register(grpcServer)
	}
	httpHandler := s.httpHandler()
	openTracing := serviceConfig.OpenTracing
	systemTracing := serviceConfig.SystemTracing
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				frame,_:=runtimei.GetCallerFrame(2)
				log.Default.Error(fmt.Sprintf(" panic: %v",r), zap.String(log.Stack, fmt.Sprintf("%s:%d (%#x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function)))
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
			if parent, ok := propagation.FromBinary(stringsi.ToBytes(r.Header.Get(httpi.HeaderTrace))); ok {
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

		if s.CustomContext != nil {
			ctx = s.CustomContext(ctx, r)
		}
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
	fmt.Printf("listening: http://%s\n", serviceConfig.Domain+serviceConfig.Port)
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
	CustomContext  CustomContext
	ConvertContext     ConvertContext
}

var signals = make(chan os.Signal, 1)
var stop = make(chan struct{}, 1)

func (s *Server) Start() {
	if initialize.InitConfig.ConfigCenter == nil {
		log.Fatal(`初始化配置失败:
	main 函数的第一行应为
	defer v2.Start(config.Conf, dao.Dao)()
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

func accessLog(ctxi pick.Context, iface, body, result string, code int) {
	ctxi.GetLogger().Logger.Info("", zap.String("interface", iface),
		zap.String("body", body),
		zap.Duration("processTime", time.Now().Sub(ctxi.GetReqTime())),
		zap.String("result", result),
		zap.String("auth", ctxi.GeToken()),
		zap.Int("status", code))
	ctxi.GetLogger().Logger.Sync()
}
