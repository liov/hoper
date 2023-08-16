package main

import (
	"github.com/actliboy/hoper/server/go/user/api"
	"github.com/actliboy/hoper/server/go/user/confdao"
	"github.com/gin-gonic/gin"
	"github.com/hopeio/zeta/initialize"
	pickgin "github.com/hopeio/zeta/pick/gin"
	"github.com/hopeio/zeta/server"
)

func main() {
	server.Start(&server.Server{
		GRPCHandle: api.GrpcRegister,

		GinHandle: func(app *gin.Engine) {
			api.GinRegister(app)
			pickgin.Register(app, confdao.Conf.Server.GenDoc, initialize.GlobalConfig.Module, confdao.Conf.Server.OpenTracing)
		},

		/*		GraphqlResolve: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/
	})
}
