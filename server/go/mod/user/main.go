package main

import (
	"context"
	"github.com/liov/hoper/server/go/lib/pick"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/server/go/lib/initialize"
	tailmon "github.com/liov/hoper/server/go/lib/tiga"
	"github.com/liov/hoper/server/go/lib/utils/net/http/gin/oauth"
	model "github.com/liov/hoper/server/go/mod/protobuf/user"
	"github.com/liov/hoper/server/go/mod/user/conf"
	"github.com/liov/hoper/server/go/mod/user/service"

	"google.golang.org/grpc"
)

func main() {
	pick.RegisterFiberService(service.GetUserService())
	app := fiber.New()
	pick.FiberWithCtx(app, true, initialize.InitConfig.Module)
	go app.Listen(":3000")
	(&tailmon.Server{
		GRPCHandle: func(gs *grpc.Server) {
			//grpc.OpenTracing = true
			model.RegisterUserServiceServer(gs, service.GetUserService())
			model.RegisterOauthServiceServer(gs, service.GetOauthService())
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {

		},
		GinHandle: func(app *gin.Engine) {
			_ = model.RegisterUserServiceHandlerServer(app, service.GetUserService())
			_ = model.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
			oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			pick.RegisterService(service.GetUserService())
			pick.Gin(app, conf.Conf.Server.GenDoc, initialize.InitConfig.Module, conf.Conf.Server.OpenTracing)
		},

		/*		GraphqlResolve: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/
	}).Start()
}
