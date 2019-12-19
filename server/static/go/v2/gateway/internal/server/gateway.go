package server

import (
	"context"
	"expvar"
	"net/http"
	"net/http/pprof"
	"os"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/http/iris/api"
	iris_log_mid "github.com/liov/hoper/go/v2/utils/http/iris/log"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/protobuf/jsonpb"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var ch = make(chan os.Signal, 1)

func SignalChan() chan os.Signal {
	return ch
}

func GateWay() http.Handler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jsonpb := &jsonpb.JSONPb{
		json.Json,
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler))
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := model.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, config.Conf.Server.GrpcService["user"], opts)
	if err != nil {
		log.Fatal(err)
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
	logger := (&log.Config{Development: config.Conf.Env == initialize.PRODUCT}).NewLogger()
	mux.Use(iris_log_mid.LogMid(logger, false))
	mux.Any("/{grpc:path}", handlerconv.FromStd(gwmux))
	mux.Get("/debug/vars", handlerconv.FromStd(expvar.Handler()))
	mux.Get("/debug/pprof/", handlerconv.FromStd(pprof.Index))
	mux.Get("/debug/pprof/cmdline", handlerconv.FromStd(pprof.Cmdline))
	mux.Get("/debug/pprof/profile", handlerconv.FromStd(pprof.Profile))
	mux.Get("/debug/pprof/symbol", handlerconv.FromStd(pprof.Symbol))
	mux.Get("/debug/pprof/trace", handlerconv.FromStd(pprof.Trace))
	api.OpenApi(mux, "../protobuf/api/")
	iris_build.Build(mux, initialize.ConfUrl)
	h2Handler := h2c.NewHandler(mux, &http2.Server{})
	server := &http.Server{Addr: config.Conf.Server.Port, Handler: h2Handler}
	go func() {
		log.Info("listening ", config.Conf.Server.Port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
	<-ch
	if err := server.Close(); err != nil {
		log.Error(err)
	}
	return mux
}
