package main

import (
	"github.com/hopeio/cherry/server"
	"github.com/liov/hoper/server/go/user/api"
)

func main() {
	server.Start(&server.Server{
		GrpcHandler: api.GrpcRegister,

		GinHandler: api.GinRegister,

		/*		GraphqlResolve: model.NewExecutableSchema(model.Config{
				Resolvers: &model.GQLServer{
					UserService:  service.GetUserService(),
					OauthService: service.GetOauthService(),
				}}),*/
	})
}
