package main

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/note/internal/config"
	"github.com/liov/hoper/go/v2/note/internal/dao"
	"github.com/liov/hoper/go/v2/note/internal/service"
	model "github.com/liov/hoper/go/v2/protobuf/note"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
	"google.golang.org/grpc"
)

func main() {
	s := server.Server{
		Conf: config.Conf,
		Dao:  dao.Dao,
		GRPCRegistr: func(g *grpc.Server) {
			model.RegisterNoteServiceServer(g, service.NoteSvc)
		},
		GatewayRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			err := model.RegisterNoteServiceHandlerServer(ctx, mux, service.NoteSvc)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	s.Start()
}
