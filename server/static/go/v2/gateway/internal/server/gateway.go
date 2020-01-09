package server

import (
	"context"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/http/iris/api"
	"github.com/liov/hoper/go/v2/utils/http/iris/gateway"
	"github.com/liov/hoper/go/v2/utils/http/iris/log"
	"github.com/liov/hoper/go/v2/utils/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
)

var ch = make(chan os.Signal, 1)

func SignalChan() chan os.Signal {
	return ch
}

func GateWay() {
	gatewayHandle := func(ctx context.Context, gwmux *runtime.ServeMux) {
		opts := []grpc.DialOption{grpc.WithInsecure()}
		err := model.RegisterUserServiceHandlerFromEndpoint(ctx, gwmux, config.Conf.Server.GrpcService["user"], opts)
		if err != nil {
			log.Fatal(err)
		}
	}
	irisHandle := func(mux *iris.Application) {
		iris_build.WithConfiguration(mux, initialize.ConfUrl)
		logger := (&log.Config{Development: config.Conf.Env == initialize.PRODUCT}).NewLogger()
		iris_log.SetLog(mux, logger, false)
		api.OpenApi(mux, "../protobuf/api/")
	}
	//openapi
	mux := iris_gateway.Http(irisHandle, gatewayHandle)
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
}
