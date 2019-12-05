package server

import (
	"context"
	"expvar"
	"mime"
	"net/http"
	"net/http/pprof"
	"os"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/gateway/internal/api"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
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

	jsonpb := &gateway.JSONPb{
		EmitDefaults: true,
		Indent:       "  ",
		OrigName:     true,
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
	mux := http.NewServeMux()
	mux.Handle("/", gwmux)
	mux.Handle("/debug/vars", expvar.Handler())
	mux.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	mux.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	mux.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	mux.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	mux.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	api.OpenApi(mux)
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
