package client

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/grpc/stats"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
	"google.golang.org/grpc"
)

func GetUserClient() user.UserServiceClient {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8090", grpc.WithInsecure(),
		grpc.WithStatsHandler(&stats.ClientHandler{}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	Connes = append(Connes, conn)
	return user.NewUserServiceClient(conn)
}
