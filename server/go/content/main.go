package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/v2/content/conf"
	"github.com/liov/hoper/v2/content/dao"
	"github.com/liov/hoper/v2/content/service"
	model "github.com/liov/hoper/v2/protobuf/content"
	"github.com/liov/hoper/v2/protobuf/user"
	"github.com/liov/hoper/v2/tiga"
	"github.com/liov/hoper/v2/tiga/initialize"
	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(conf.Config, dao.Dao)()

	s := tiga.Server{
		GRPCHandle: func(gs *grpc.Server) {
			model.RegisterMomentServiceServer(gs, service.GetMomentService())
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			_ = model.RegisterMomentServiceHandlerServer(ctx, mux, service.GetMomentService())

		},
		CustomContext:  user.CtxWithRequest,
		ConvertContext: user.Authorization,
	}
	s.Start()
}
