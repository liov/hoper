package server

import (
	"context"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/service"
	"github.com/liov/hoper/go/v2/utils/http/gateway"
	"github.com/liov/hoper/go/v2/utils/log"
)

func Http() http.Handler {
	gatewayHandle := func(ctx context.Context, gwmux *runtime.ServeMux) {
		runtime.WithForwardResponseOption(hook)(gwmux)
		err := model.RegisterUserServiceHandlerServer(ctx, gwmux, service.UserSvc)
		if err != nil {
			log.Fatal(err)
		}
	}
	mux := gateway.Http(nil, gatewayHandle, config.Conf.Env)
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
