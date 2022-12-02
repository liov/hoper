package main

import (
	"context"

	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga"
	"github.com/actliboy/hoper/server/go/mod/content/conf"
	"github.com/actliboy/hoper/server/go/mod/content/dao"
	"github.com/actliboy/hoper/server/go/mod/content/service"
	model "github.com/actliboy/hoper/server/go/mod/protobuf/content"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
