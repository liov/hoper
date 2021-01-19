package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/user/service"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/oauth"
	igrpc "github.com/liov/hoper/go/v2/utils/net/http/grpc"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/liov/hoper/go/v2/utils/net/http/tailmon"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(conf.Conf, dao.Dao)()
	pick.RegisterService(service.GetUserService())
	(&tailmon.Server{
		//为了可以自定义中间件
		GRPCServer: func() *grpc.Server {
			gs := igrpc.DefaultGRPCServer(nil,nil)
			model.RegisterUserServiceServer(gs, service.GetUserService())
			model.RegisterOauthServiceServer(gs,service.GetOauthService())
			return gs
		}(),
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			if err := model.RegisterUserServiceHandlerServer(ctx, mux, service.GetUserService()); err != nil {
				log.Fatal(err)
			}
			if err := model.RegisterOauthServiceHandlerServer(ctx, mux, service.GetOauthService()); err != nil {
				log.Fatal(err)
			}
		},
		GinHandle: func(app *gin.Engine) {
			oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			pick.Gin(app,true,initialize.InitConfig.Module)
		},
	}).Start(service.CtxWithRequest)
}