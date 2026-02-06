package global

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
	UserClient = sync.OnceValue(func() user.UserServiceClient {
		// Set up a connection to the server.
		conn, err := grpcx.NewClient("127.0.0.1:8080",
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return user.NewUserServiceClient(conn)
	})
	UploadClient = sync.OnceValue(func() file.FileServiceClient {
		// Set up a connection to the server.
		conn, err := grpcx.NewClient("127.0.0.1:8080",
			grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		)

		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return file.NewFileServiceClient(conn)
	})
)
