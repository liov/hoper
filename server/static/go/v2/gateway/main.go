package main

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	note "github.com/liov/hoper/go/v2/protobuf/note"
	user "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/server"
	"google.golang.org/grpc"
)

func main() {
	s := server.Server{
		Conf:        config.Conf,
		Dao:         nil,
		GRPCRegistr: nil,
		HTTPRegistr: func(ctx context.Context, mux *runtime.ServeMux) {
			opts := []grpc.DialOption{grpc.WithInsecure()}
			err := user.RegisterUserServiceHandlerFromEndpoint(ctx, mux, config.Conf.Customize.GrpcService["user"], opts)
			if err != nil {
				log.Fatal(err)
			}
			err = note.RegisterNoteServiceHandlerFromEndpoint(ctx, mux, config.Conf.Customize.GrpcService["note"], opts)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	s.Start()
}
