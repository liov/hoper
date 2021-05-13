package client

import (
	model "github.com/liov/hoper/v2/protobuf/user"
	"github.com/liov/hoper/v2/utils/log"
	"github.com/liov/hoper/v2/utils/net/http/grpc/stats"
	"google.golang.org/grpc"
)

func GetUserClient() (model.UserServiceClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8090", grpc.WithInsecure(),
		grpc.WithStatsHandler(&stats.ClientHandler{}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return model.NewUserServiceClient(conn), conn
}
