package rpc

import (
	grpci "github.com/hopeio/lemon/utils/net/http/grpc"
	"github.com/liovx/hoper/server/go/protobuf/user"
	"log"
)

var (
	UserClient user.UserServiceClient
)

func init() {
	UserClient = GetUserClient()
}

func GetUserClient() user.UserServiceClient {

	// Set up a connection to the server.
	conn, err := grpci.GetTlsClient("grpc.hoper.xyz:443")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return user.NewUserServiceClient(conn)
}
