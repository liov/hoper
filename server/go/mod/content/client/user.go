package client

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	grpci "github.com/liov/hoper/server/go/lib/utils/net/http/grpc"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
)

func GetUserClient() user.UserServiceClient {
	// Set up a connection to the server.
	conn, err := grpci.GetDefaultClient("localhost:8090")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	Connes = append(Connes, conn)
	return user.NewUserServiceClient(conn)
}
