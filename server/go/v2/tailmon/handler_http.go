package tailmon

import (
	"bytes"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	"github.com/liov/hoper/go/v2/tailmon/initialize"

	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	gin_build "github.com/liov/hoper/go/v2/utils/net/http/gin"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/gateway"

	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func (s *Server) httpHandler(conf *initialize.ServerConfig) http.HandlerFunc {
	// 默认使用gin
	ginServer := gin_build.Http(conf.Gin, s.GinHandle)

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
	var excludes = []string{"/api/v1/upload", "/api/v1/multiUpload"}
	var includes = []string{"/api"}
	return func(w http.ResponseWriter, r *http.Request) {
		// 暂时解决方法，三个路由
		if h, p := http.DefaultServeMux.Handler(r); p != "" {
			h.ServeHTTP(w, r)
		}
		if !stringsi.HasPrefixes(r.RequestURI, includes) || stringsi.HasPrefixes(r.RequestURI, excludes) {
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

		accessLog(contexti.CtxFromContext(r.Context()), r.RequestURI, r.Method,
			stringsi.ToString(body), stringsi.ToString(recorder.Body.Bytes()),
			recorder.Code)
	}
}

func accessLog(ctxi *contexti.Ctx, iface, method, body, result string, code int) {
	// log 里time now 浪费性能
	if ce := ctxi.GetLogger().Logger.Check(zap.InfoLevel, "access"); ce != nil {
		ce.Write(zap.String("interface", iface),
			zap.String("method", method),
			zap.String("body", body),
			// 性能
			zap.Duration("processTime", ce.Time.Sub(ctxi.GetReqAt().Time)),
			zap.String("result", result),
			zap.String("auth", ctxi.GeToken()),
			zap.Int("status", code))
	}
}
