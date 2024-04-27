package rpc

import (
	grpci "github.com/hopeio/cherry/utils/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"sync"
)

var (
	lock       sync.Mutex
	userClient user.UserServiceClient
)

func UserClient() user.UserServiceClient {
	if userClient != nil {
		return userClient
	}
	lock.Lock()

	// Set up a connection to the server.
	conn, err := grpci.GetTlsClient("grpc.hoper.xyz:443", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	userClient = user.NewUserServiceClient(conn)
	lock.Unlock()
	return userClient
}
