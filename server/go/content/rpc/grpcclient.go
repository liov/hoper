package rpc

import (
	"github.com/hopeio/cherry/utils/log"
	grpci "github.com/hopeio/cherry/utils/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/upload"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"sync"
)

var (
	userClient   user.UserServiceClient
	uploadClient upload.UploadServiceClient
)

func UserClient() user.UserServiceClient {
	if userClient != nil {
		return userClient
	}
	userClient = sync.OnceValue[user.UserServiceClient](func() user.UserServiceClient {
		// Set up a connection to the server.
		conn, err := grpci.GetDefaultClient("localhost:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return user.NewUserServiceClient(conn)
	})()
	return userClient
}

func UploadClient() upload.UploadServiceClient {
	if uploadClient != nil {
		return uploadClient
	}
	uploadClient = sync.OnceValue[upload.UploadServiceClient](func() upload.UploadServiceClient {
		// Set up a connection to the server.
		conn, err := grpci.GetDefaultClient("localhost:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return upload.NewUploadServiceClient(conn)
	})()
	return uploadClient
}
