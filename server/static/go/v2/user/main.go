package main

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/user/internal/service"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/filter"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer v2.Start(config.Conf, dao.Dao)()
	s := server.Server{
		//为了可以自定义中间件
		GRPCServer: func() *grpc.Server {
			gs := grpc.NewServer(
				//filter应该在最前
				grpc.UnaryInterceptor(
					grpc_middleware.ChainUnaryServer(
						filter.UnaryServerInterceptor()...,
					)),
				grpc.StreamInterceptor(
					grpc_middleware.ChainStreamServer(
						filter.StreamServerInterceptor()...,
					)),
			)
			model.RegisterUserServiceServer(gs, service.GetUserService())
			return gs
		}(),
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			if err := model.RegisterUserServiceHandlerServer(ctx, mux, service.GetUserService()); err != nil {
				log.Fatal(err)
			}
		},
		PickHandle: func(app *pick.EasyRouter) {
			/*			oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
						app.Get("/oauth/login", func(ctx context2.Context) {
							ctx.ServeFile("./static/login.html", false)
						})*/
		},
		GraphqlResolve: model.NewExecutableSchema(model.Config{
			Resolvers: &model.GQLServer{
				UserService:  service.GetUserService(),
				OauthService: service.GetOauthService(),
			}}),
	}
	s.Start()
}
