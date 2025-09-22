package data

import (
	grpcx "github.com/hopeio/gox/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/user"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"log"
	"sync"
)

var (
	UserClient = sync.OnceValue(func() user.UserServiceClient {
		conn, err := grpcx.NewTLSClient("grpc.hoper.xyz:443", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		return user.NewUserServiceClient(conn)
	})
)
