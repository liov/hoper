package server

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/gateway/response"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/service"
	iris_build "github.com/liov/hoper/go/v2/utils/http/iris"
	"github.com/liov/hoper/go/v2/utils/http/iris/api"
	"github.com/liov/hoper/go/v2/utils/http/iris/gateway"
	iris_log "github.com/liov/hoper/go/v2/utils/http/iris/log"
	"github.com/liov/hoper/go/v2/utils/log"
)

func Http() http.Handler {
	gatewayHandle := func(ctx context.Context, gwmux *runtime.ServeMux) {
		runtime.WithForwardResponseOption(response.UserHook(config.Conf.Server.TokenMaxAge))(gwmux)
		err := model.RegisterUserServiceHandlerServer(ctx, gwmux, service.UserSvc)
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
	mux := iris_gateway.Http(irisHandle, gatewayHandle)
	return mux
}
