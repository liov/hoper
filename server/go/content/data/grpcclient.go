package data

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
	UserClient   = sync.OnceValue[user.UserServiceClient](func() user.UserServiceClient {
		// Set up a connection to the server.
		conn, err := grpci.NewClient("127.0.0.1:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return user.NewUserServiceClient(conn)
	})
	UploadClient = sync.OnceValue[upload.UploadServiceClient](func() upload.UploadServiceClient {
		// Set up a connection to the server.
		conn, err := grpci.NewClient("127.0.0.1:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return upload.NewUploadServiceClient(conn)
	})
)
