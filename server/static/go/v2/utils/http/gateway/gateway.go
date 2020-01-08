package gateway

import (
	"context"
	"expvar"
	"net/http"
	"net/http/pprof"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/initialize"

	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/http/iris/api"
	iris_log_mid "github.com/liov/hoper/go/v2/utils/http/iris/log"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/protobuf/jsonpb"
	"google.golang.org/grpc/metadata"
)

func Http(irisHandle func(*iris.Application), gatewayHandle func(context.Context, *runtime.ServeMux), env string) http.Handler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jsonpb := &jsonpb.JSONPb{
		json.Json,
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			return map[string][]string{
				"ip":        {request.RemoteAddr},
				"UserAgent": {request.UserAgent()},
			}
		}),
		runtime.WithIncomingHeaderMatcher(func(s2 string) (s string, b bool) {
			if s2 == "Cookie" {
				return s2, true
			} else {
				return s2, false
			}
		}))
	if gatewayHandle != nil {
		gatewayHandle(ctx, gwmux)
	}
	//openapi
	mux := iris.New()

	iris.RegisterOnInterrupt(func() {
		timeout := 5 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		//关闭所有主机
		mux.Shutdown(ctx)
	})
	logger := (&log.Config{Development: env == initialize.PRODUCT}).NewLogger()
	mux.Use(iris_log_mid.LogMid(logger, false))
	mux.Any("/{grpc:path}", handlerconv.FromStd(gwmux))
	mux.Get("/debug/vars", handlerconv.FromStd(expvar.Handler()))
	mux.Get("/debug/pprof/", handlerconv.FromStd(pprof.Index))
	mux.Get("/debug/pprof/cmdline", handlerconv.FromStd(pprof.Cmdline))
	mux.Get("/debug/pprof/profile", handlerconv.FromStd(pprof.Profile))
	mux.Get("/debug/pprof/symbol", handlerconv.FromStd(pprof.Symbol))
	mux.Get("/debug/pprof/trace", handlerconv.FromStd(pprof.Trace))
	api.OpenApi(mux, "../protobuf/api/")
	if irisHandle != nil {
		irisHandle(mux)
	}
	iris_build.Build(mux, initialize.ConfUrl)
	return mux
}
