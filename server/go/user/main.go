package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hopeio/tiga/server"
	"github.com/liov/hoper/server/go/user/api"
)

func main() {
	server.Start(&server.Server{
		GRPCHandle: api.GrpcRegister,

		GinHandle: func(app *gin.Engine) {
			api.GinRegister(app)
		},

		/*		GraphqlResolve: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/
	})
}
