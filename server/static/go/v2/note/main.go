package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/note/internal/config"
	"github.com/liov/hoper/go/v2/note/internal/dao"
	"github.com/liov/hoper/go/v2/note/internal/service"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/server"
	"google.golang.org/grpc"
)

func main() {
	s := server.Server{
		Conf: config.Conf,
		Dao:  dao.Dao,
		GRPCRegistr: func(g *grpc.Server) {
			model.RegisterUserServiceServer(g, service.NoteSvc)
		},
		HTTPRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			err := model.RegisterUserServiceHandlerServer(ctx, mux, service.NoteSvc)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	s.Start()
}
