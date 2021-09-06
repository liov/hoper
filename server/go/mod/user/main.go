package main

import (
	"context"
	"github.com/liov/hoper/server/go/lib/tiga/pick"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	tailmon "github.com/liov/hoper/server/go/lib/tiga"
	"github.com/liov/hoper/server/go/lib/tiga/initialize"
	"github.com/liov/hoper/server/go/lib/utils/net/http/gin/oauth"
	model "github.com/liov/hoper/server/go/mod/protobuf/user"
	"github.com/liov/hoper/server/go/mod/user/conf"
	"github.com/liov/hoper/server/go/mod/user/dao"
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
			_ = model.RegisterUserServiceHandlerServer(ctx, mux, service.GetUserService())
			_ = model.RegisterOauthServiceHandlerServer(ctx, mux, service.GetOauthService())
		},
		GinHandle: func(app *gin.Engine) {
			oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
			app.StaticFS("/oauth/login", http.Dir("./static/login.html"))
			pick.RegisterService(service.GetUserService())
			pick.Gin(app, true, initialize.InitConfig.Module)
		},

		/*		GraphqlResolve: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/
	}).Start()
}