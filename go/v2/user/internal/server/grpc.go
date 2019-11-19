package server

import (
	"context"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/liov/hoper/go/v2/protobuf/response"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func filter(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("%v panic: %v", info, r)
		}
		if err != nil {
			resp = &response.AnyReply{
				Code: 1000,
				Data: nil,
				Msg:  err.Error(),
			}
			s,_ := status.FromError(err)
			s,_ = s.WithDetails(&model.User{FollowCount: 10})
			err = s.Err()
		}
	}()

	return handler(ctx, req)
}

func Grpc() *grpc.Server {
	s := grpc.NewServer(
		//filter应该在最前
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(filter,grpc_validator.UnaryServerInterceptor())),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(grpc_validator.StreamServerInterceptor())),
	)
	model.RegisterUserServiceServer(s, userService)
	// Register reflection service on gRPC server.
	reflection.Register(s)
	return s
}
