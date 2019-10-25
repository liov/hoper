package client

import (
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc"
)

func GetUserClient() user.UserServiceClient {
	// Set up a connection to the server.
	conn, err := grpc.Dial("user.hoper.xyz:8080", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	return user.NewUserServiceClient(conn)
}
