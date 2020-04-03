package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/user/internal/service"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
	"google.golang.org/grpc"
)

func main() {
	s := server.Server{
		Conf: config.Conf,
		Dao:  dao.Dao,
		GRPCRegistr: func(g *grpc.Server) {
			model.RegisterUserServiceServer(g, service.UserSvc)
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			err := model.RegisterUserServiceHandlerServer(ctx, mux, service.UserSvc)
			if err != nil {
				log.Fatal(err)
			}
		},
		GraphqlResolve: service.UserSvc,
	}
	s.Start()
}
