package server

import (
	"context"
	"net/http"
	"time"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/service"
	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/http/iris/api"
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
	err := model.RegisterUserServiceHandlerServer(ctx, gwmux, service.UserSvc)
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
	api.OpenApi(mux, "../protobuf/api/")
	iris_build.Build(mux, initialize.ConfUrl)
	return mux
}
