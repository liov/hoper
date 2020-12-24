package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/note/conf"
	"github.com/liov/hoper/go/v2/note/dao"
	"github.com/liov/hoper/go/v2/note/service"
	model "github.com/liov/hoper/go/v2/protobuf/note"
	"github.com/liov/hoper/go/v2/utils/log"
	igrpc "github.com/liov/hoper/go/v2/utils/net/http/grpc"
	"github.com/liov/hoper/go/v2/utils/net/http/tailmon"
	"google.golang.org/grpc"
)

func main() {
	defer initialize.Start(conf.Config, dao.Dao)()

	s := tailmon.Server{
		GRPCServer: func() *grpc.Server {
			gs := igrpc.DefaultGRPCServer(nil,nil)
			model.RegisterNoteServiceServer(gs, service.NoteSvc)
			return gs
		}(),
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			err := model.RegisterNoteServiceHandlerServer(ctx, mux, service.NoteSvc)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	s.Start()
}
