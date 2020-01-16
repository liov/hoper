package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/service"
	"github.com/liov/hoper/go/v2/utils/grpc/filter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func Grpc() *grpc.Server {
	s := grpc.NewServer(
		//filter应该在最前
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(filter.CommonUnaryServerInterceptor()...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(filter.CommonStreamServerInterceptor()...)),
	)
	model.RegisterUserServiceServer(s, service.UserSvc)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	return s
}
