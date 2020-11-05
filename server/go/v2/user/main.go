package main

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/kataras/iris/v12"
	context2 "github.com/kataras/iris/v12/context"
	"github.com/liov/hoper/go/v2/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/user/service"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/filter"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/oauth"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
	"google.golang.org/grpc"
)

func main() {
	//配置初始化应该在第一位
	defer initialize.Start(conf.Config, dao.Dao)()
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
		IrisHandle: func(app *iris.Application) {
			oauth.RegisterOauthServiceHandlerServer(app, service.GetOauthService())
			app.Get("/oauth/login", func(ctx *context2.Context) {
				ctx.ServeFile("./static/login.html")
			})
		},
		GraphqlResolve: model.NewExecutableSchema(model.Config{
			Resolvers: &model.GQLServer{
				UserService:  service.GetUserService(),
				OauthService: service.GetOauthService(),
			}}),
	}
	s.Start()
}
