package main

import (
	"github.com/actliboy/hoper/server/go/user/api"
	"github.com/actliboy/hoper/server/go/user/confdao"
	"github.com/actliboy/hoper/server/go/user/service"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/zeta/initialize"
	"github.com/hopeio/zeta/pick"
	"github.com/hopeio/zeta/server"
)

func main() {
	pick.RegisterFiberService(service.GetUserService())
	app := fiber.New()
	pick.FiberWithCtx(app, true, initialize.GlobalConfig.Module)
	go app.Listen(":3000")
	server.Start(&server.Server{
		GRPCHandle: api.GrpcRegister,

		GinHandle: func(app *gin.Engine) {
			api.GinRegister(app)
			pick.Gin(app, confdao.Conf.Server.GenDoc, initialize.GlobalConfig.Module, confdao.Conf.Server.OpenTracing)
		},

		/*		GraphqlResolve: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/
	})
}
