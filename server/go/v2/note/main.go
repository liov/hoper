package main

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	"github.com/liov/hoper/go/v2/note/internal/config"
	"github.com/liov/hoper/go/v2/note/internal/dao"
	"github.com/liov/hoper/go/v2/note/internal/service"
	model "github.com/liov/hoper/go/v2/protobuf/note"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/filter"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
	"google.golang.org/grpc"
)

func main() {
	defer v2.Start(config.Conf, dao.Dao)()

	s := server.Server{
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