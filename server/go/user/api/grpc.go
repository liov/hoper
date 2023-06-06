package api

import (
	"github.com/actliboy/hoper/server/go/protobuf/user"
	userService "github.com/actliboy/hoper/server/go/user/service"
	"google.golang.org/grpc"
)

func GrpcRegister(gs *grpc.Server) {
	user.RegisterUserServiceServer(gs, userService.GetUserService())
	user.RegisterOauthServiceServer(gs, userService.GetOauthService())
}
