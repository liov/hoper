package client

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	model "github.com/liov/hoper/server/go/mod/protobuf/user"
	"google.golang.org/grpc"
)

func GetUserClient() (model.UserServiceClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8090", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return model.NewUserServiceClient(conn), conn
}
