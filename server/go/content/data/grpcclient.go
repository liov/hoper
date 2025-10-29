package data

import (
	"sync"

	"github.com/hopeio/gox/log"
	grpcx "github.com/hopeio/gox/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/file"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var (
	userClient user.UserServiceClient
	fileClient file.FileServiceClient
	UserClient = sync.OnceValue[user.UserServiceClient](func() user.UserServiceClient {
		// Set up a connection to the server.
		conn, err := grpcx.NewClient("127.0.0.1:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return user.NewUserServiceClient(conn)
	})
	UploadClient = sync.OnceValue[file.FileServiceClient](func() file.FileServiceClient {
		// Set up a connection to the server.
		conn, err := grpcx.NewClient("127.0.0.1:8090", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return file.NewFileServiceClient(conn)
	})
)
