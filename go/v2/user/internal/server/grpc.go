package server

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/service"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func filter(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("%v panic: %v",info,r)
		}
	}()

	return handler(ctx, req)
}

func Grpc() *grpc.Server{
	s := grpc.NewServer(
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(filter,grpc_validator.UnaryServerInterceptor())),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(grpc_validator.StreamServerInterceptor())),
		)
	model.RegisterUserServiceServer(s, &service.UserService{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	return  s
}

