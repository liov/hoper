package server

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/core/handlerconv"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/service"
	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/http/iris/api"
	iris_log_mid "github.com/liov/hoper/go/v2/utils/http/iris/log"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/protobuf/jsonpb"
)

func Http() http.Handler {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	jsonpb := &jsonpb.JSONPb{
		json.Json,
	}
	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
		runtime.WithForwardResponseOption(hook))
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
	logger := (&log.Config{Development: config.Conf.Env == initialize.PRODUCT}).NewLogger()
	mux.Use(iris_log_mid.LogMid(logger, false))
	//mux.Use(handlerconv.FromStd(middle.HttpAuth))
	mux.Any("/{grpc:path}", handlerconv.FromStd(gwmux))
	api.OpenApi(mux, "../protobuf/api/")
	iris_build.Build(mux, initialize.ConfUrl)
	return mux
}

func hook(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	switch rep := message.(type) {
	case *model.LoginRep:
		if rep.Details != nil {
			http.SetCookie(writer, &http.Cookie{
				Name:  "token",
				Value: rep.Details.Token,
				Path:  "/",
				//Domain:   "hoper.xyz",
				Expires:  time.Now().Add(time.Duration(config.Conf.Server.TokenMaxAge) * time.Second),
				MaxAge:   int(time.Duration(config.Conf.Server.TokenMaxAge) * time.Second),
				Secure:   false,
				HttpOnly: true,
			})
		}
	case *model.LogoutRep:
		http.SetCookie(writer, &http.Cookie{
			Name:  "token",
			Value: "del",
			Path:  "/",
			//Domain:   "hoper.xyz",
			Expires:  time.Now().Add(-1),
			MaxAge:   -1,
			Secure:   false,
			HttpOnly: true,
		})
	}
	return nil
}
