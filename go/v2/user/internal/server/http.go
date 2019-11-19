package server

import (
	"context"
	"net/http"

	"github.com/gogo/gateway"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	model "github.com/liov/hoper/go/v2/protobuf/user"
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
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),)
	err := model.RegisterUserServiceHandlerServer(ctx,gwmux,userService)
	if err != nil {
		log.Fatal(err)
	}
	//openapi
	mux := http.NewServeMux()
	mux.Handle("/",gwmux)
	mux.Handle("/open-api/",http.StripPrefix("/open-api/", http.FileServer(http.Dir("./api"))))
	return mux
}
