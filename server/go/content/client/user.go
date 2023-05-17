package client

import (
	"github.com/actliboy/hoper/server/go/protobuf/user"
	"github.com/hopeio/pandora/utils/log"
	grpci "github.com/hopeio/pandora/utils/net/http/grpc"
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
