package global

import (
	"sync"

	"github.com/hopeio/gox/log"
	grpcx "github.com/hopeio/gox/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/file"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/stats/opentelemetry"
)

var (
	option     = otelgrpc.WithPropagators(
		propagation.NewCompositeTextMapPropagator(
			opentelemetry.GRPCTraceBinPropagator{}, propagation.Baggage{},
		))
	UserClient = sync.OnceValue(func() user.UserServiceClient {
		// Set up a connection to the server.
		conn, err := grpcx.NewClient("127.0.0.1:8080",
			grpc.WithStatsHandler(otelgrpc.NewClientHandler(option)))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return user.NewUserServiceClient(conn)
	})
	UploadClient = sync.OnceValue(func() file.FileServiceClient {
		// Set up a connection to the server.
		conn, err := grpcx.NewClient("127.0.0.1:8080",
			grpc.WithStatsHandler(otelgrpc.NewClientHandler(option)),
		)
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return file.NewFileServiceClient(conn)
	})
)
