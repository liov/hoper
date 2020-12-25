package tailmon

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/signal"
	"runtime/debug"
	"strings"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/log"
	gin_build "github.com/liov/hoper/go/v2/utils/net/http/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/liov/hoper/go/v2/utils/strings2"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Server) httpHandler() http.HandlerFunc {
	var graphqlServer http.Handler
	if s.GraphqlResolve != nil {
		graphqlServer = handler.NewDefaultServer(s.GraphqlResolve)
	}
	var gatewayServer http.Handler
	if s.GatewayRegistr != nil {
		gatewayServer = gateway.Gateway(s.GatewayRegistr)
	}
	var ginServer *gin.Engine
	if s.GinHandle != nil {
		ginServer = gin_build.Http(initialize.InitConfig.ConfUrl, "../protobuf/api/", s.GinHandle)
	}
	var pickServer *pick.Router
	if s.PickHandle !=nil{
		pickServer = pick.New(false, initialize.InitConfig.Module)
	}
	return func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		/*		var result bytes.Buffer
				rsp := io.MultiWriter(w, &result)*/
		recorder := httptest.NewRecorder()
		// 根据header判断走哪个runtime 而不是gin统一代理
		if gatewayServer != nil && r.Header.Get("Runtime") == "gateway" {
			gatewayServer.ServeHTTP(recorder, r)
			return
		}else if graphqlServer != nil && r.Header.Get("Runtime") == "graphql" {
			graphqlServer.ServeHTTP(recorder, r)
			return
		}else {
			ginServer.ServeHTTP(recorder, r)
		}

		if recorder.Code == http.StatusNotFound && pickServer !=nil {
			recorder = httptest.NewRecorder()
			pickServer.ServeHTTP(recorder,r)
		}

		// 从 recorder 中提取记录下来的 Response Header，设置为 ResponseWriter 的 Header
		for key, value := range recorder.Header() {
			for _, val := range value {
				w.Header().Set(key, val)
			}
		}

		// 提取 recorder 中记录的状态码，写入到 ResponseWriter 中
		w.WriteHeader(recorder.Code)
		if recorder.Body != nil {
			// 将 recorder 记录的 Response Body 写入到 ResponseWriter 中，客户端收到响应报文体
			w.Write(recorder.Body.Bytes())
		}

		AccessLog(strings2.ToSting(recorder.Body.Bytes()), r, "Cookie", now)
	}
}

func (s *Server) Serve() {
	if s.GRPCServer != nil {
		reflection.Register(s.GRPCServer)
	}
	httpHandler := s.httpHandler()
	handle := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if r := recover(); r != nil {
				log.CallTwo.With(zap.String("stack", strings2.ToSting(debug.Stack()))).Error(" panic: ", r)
				w.Write(errorcode.SysErr)
			}
		}()
		// 请求UUID，链路跟踪用
		r.Header.Set("UUID", uuid.New().String())

		if r.ProtoMajor == 2 && s.GRPCServer != nil && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			s.GRPCServer.ServeHTTP(w, r) // gRPC Server
		} else {
			httpHandler(w, r)
		}
	})
	h2Handler := h2c.NewHandler(handle, &http2.Server{})
	//反射从配置中取port
	serviceConfig := initialize.InitConfig.GetServiceConfig()
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
	log.Debugf("listening%v", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Server struct {
	GRPCServer     *grpc.Server
	GatewayRegistr gateway.GatewayHandle
	GinHandle      func(engine *gin.Engine)
	PickHandle     func(engine *pick.Router)
	GraphqlResolve graphql.ExecutableSchema
}

var close = make(chan os.Signal, 1)
var stop = make(chan struct{}, 1)

func (s *Server) Start() {
	if initialize.InitConfig.ConfigCenter == nil {
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

func AccessLog(resp string, r *http.Request, authKey string, start time.Time) {
	//参数处理
	var param string
	if r.Method == "GET" {
		param = r.URL.RawQuery
	} else {
		rawData, _ := ioutil.ReadAll(r.Body)
		param = strings2.ToSting(rawData)
	}
	//头信息处理
	//获取请求头的报文信息
	authHeader := r.Header.Get(authKey)

	log.Default.With(
		zap.String("interface", r.RequestURI),
		zap.String("param", param),
		zap.Duration("processTime", time.Now().Sub(start)),
		zap.String("result", resp),
		zap.String("other", authHeader),
		zap.String("source", initialize.InitConfig.Module),
	).Info()
}
