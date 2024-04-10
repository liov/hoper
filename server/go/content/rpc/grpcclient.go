package rpc

import (
	"github.com/hopeio/tiga/utils/log"
	grpci "github.com/hopeio/tiga/utils/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/upload"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var (
	UserClient   user.UserServiceClient
	UploadClient upload.UploadServiceClient
)

func init() {
	UserClient = GetUserClient()
	UploadClient = GetUploadClient()
}

func GetUserClient() user.UserServiceClient {
	// Set up a connection to the server.
	conn, err := grpci.GetDefaultClient("localhost:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return user.NewUserServiceClient(conn)
}

func GetUploadClient() upload.UploadServiceClient {
	// Set up a connection to the server.
	conn, err := grpci.GetDefaultClient("localhost:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return upload.NewUploadServiceClient(conn)
}
