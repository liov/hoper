package main

import (
	"context"
	"log"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/gateway/internal/config"
	note "github.com/liov/hoper/go/v2/protobuf/note"
	user "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/net/http/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
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
			err = note.RegisterNoteServiceHandlerFromEndpoint(ctx, mux, "localhost:8001", opts)
			if err != nil {
				log.Fatal(err)
			}
		},
	}
	s.Start()
}

func service() {
	svc := map[string]func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error){
		"user": user.RegisterUserServiceHandlerFromEndpoint,
		"note": note.RegisterNoteServiceHandlerFromEndpoint,
	}
	log.Println(svc)
}

func RegisterNoteServiceHandlerFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption, f func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error) (err error) {
	conn, err := grpc.Dial(endpoint, opts...)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
			return
		}
		go func() {
			<-ctx.Done()
			if cerr := conn.Close(); cerr != nil {
				grpclog.Infof("Failed to close conn to %s: %v", endpoint, cerr)
			}
		}()
	}()

	return f(ctx, mux, conn)
}
