package main

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/user/internal/service"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/filter"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
	"google.golang.org/grpc"
)

func main() {
	s := server.Server{
		Conf: config.Conf,
		Dao:  dao.Dao,
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
			model.RegisterUserServiceServer(gs, service.UserSvc)
			return gs
		}(),
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			err := model.RegisterUserServiceHandlerServer(ctx, mux, service.UserSvc)
			if err != nil {
				log.Fatal(err)
			}
		},
		GraphqlResolve: model.NewExecutableSchema(model.Config{
			Resolvers: &model.UserServiceGQLServer{Service: service.UserSvc}}),
	}
	s.Start()
}
