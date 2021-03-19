package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/content/conf"
	"github.com/liov/hoper/go/v2/content/dao"
	"github.com/liov/hoper/go/v2/content/service"
	"github.com/liov/hoper/go/v2/tailmon"
	"github.com/liov/hoper/go/v2/tailmon/initialize"
	model "github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(conf.Config, dao.Dao)()

	s := tailmon.Server{
		GRPCHandle: func(gs *grpc.Server)  {
			model.RegisterMomentServiceServer(gs, service.GetMomentService())
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			_= model.RegisterMomentServiceHandlerServer(ctx, mux, service.GetMomentService())

		},
		CustomContext: user.CtxWithRequest,
		ConvertContext: user.Authorization,
	}
	s.Start()
}
