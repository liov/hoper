package main

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	contentervice "github.com/liov/hoper/go/v2/content/service"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	userservice "github.com/liov/hoper/go/v2/user/service"
	"github.com/liov/hoper/go/v2/utils/net/http/gin/oauth"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/liov/hoper/go/v2/utils/net/http/tailmon"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(conf.Conf, dao.Dao)()
	pick.RegisterService(userservice.GetUserService(), contentervice.GetMomentService())
	(&tailmon.Server{
		//为了可以自定义中间件
		GRPCHandle: func(gs *grpc.Server)  {
			user.RegisterUserServiceServer(gs, userservice.GetUserService())
			user.RegisterOauthServiceServer(gs, userservice.GetOauthService())
			content.RegisterMomentServiceServer(gs, contentervice.GetMomentService())
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			_ = user.RegisterUserServiceHandlerServer(ctx, mux, userservice.GetUserService())
			_ = user.RegisterOauthServiceHandlerServer(ctx, mux, userservice.GetOauthService())
			_ = content.RegisterMomentServiceHandlerServer(ctx, mux, contentervice.GetMomentService())

		},
		GinHandle: func(app *gin.Engine) {
			oauth.RegisterOauthServiceHandlerServer(app, userservice.GetOauthService())
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			pick.Gin(app, user.ConvertContext, true, initialize.InitConfig.Module)
		},
		CustomContext: user.CtxWithRequest,
		Authorization: user.Authorization,
	}).Start()
}
