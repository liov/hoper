package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hopeio/pandora/initialize"
	"github.com/hopeio/pandora/tiga"
	"github.com/liov/hoper/server/go/mod/content/conf"
	"github.com/liov/hoper/server/go/mod/content/dao"
	"github.com/liov/hoper/server/go/mod/content/service"
	model "github.com/liov/hoper/server/go/mod/protobuf/content"
	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(conf.Conf, dao.Dao)()

	s := tiga.Server{
		Config: conf.Conf.Server.Origin(),
		GRPCHandle: func(gs *grpc.Server) {
			model.RegisterMomentServiceServer(gs, service.GetMomentService())
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			_ = model.RegisterMomentServiceHandlerServer(ctx, mux, service.GetMomentService())

		},
	}
	s.Start()
}
