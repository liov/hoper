package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/api"
	"github.com/liov/hoper/go/v2/utils/log"
)

func Http() http.Handler {
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
	err := model.RegisterUserServiceHandlerServer(ctx, gwmux, userService)
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
	mux.Any("/{grpc:path}", handlerconv.FromStd(gwmux))
	api.OpenApi(mux)
	if err := mux.Build(); err != nil {
		log.Fatal(err)
	}
	//将来这块用配置中心
	mux.Configure(iris.WithConfiguration(iris.YAML("./config/iris.yml")))
	return mux
}
