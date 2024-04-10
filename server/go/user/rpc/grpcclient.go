package rpc

import (
	grpci "github.com/hopeio/tiga/utils/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
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
	conn, err := grpci.GetTlsClient("grpc.hoper.xyz:443", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return user.NewUserServiceClient(conn)
}
