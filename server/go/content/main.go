package main

import (
	"context"

	"github.com/actliboy/hoper/server/go/content/confdao"
	"github.com/actliboy/hoper/server/go/content/service"
	model "github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hopeio/pandora/initialize"
	"github.com/hopeio/pandora/server"
	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(confdao.Conf, confdao.Dao)()

	s := server.Server{
		Config: confdao.Conf.Server.Origin(),
		GRPCHandle: func(gs *grpc.Server) {
			model.RegisterMomentServiceServer(gs, service.GetMomentService())
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			_ = model.RegisterMomentServiceHandlerServer(ctx, mux, service.GetMomentService())

		},
	}
	s.Start()
}
