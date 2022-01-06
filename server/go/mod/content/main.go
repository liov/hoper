package main

import (
	"context"

	"github.com/actliboy/hoper/server/go/lib/tiga"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/mod/content/conf"
	"github.com/actliboy/hoper/server/go/mod/content/dao"
	"github.com/actliboy/hoper/server/go/mod/content/service"
	model "github.com/actliboy/hoper/server/go/mod/protobuf/content"
	"github.com/actliboy/hoper/server/go/mod/protobuf/user"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
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
