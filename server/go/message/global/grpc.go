package global

import (
	"sync"

	"github.com/hopeio/gox/log"
	grpcx "github.com/hopeio/gox/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/message"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
)

var MessageClient = sync.OnceValue(func() message.MessageClient {
	// Set up a connection to the server.
	conn, err := grpcx.NewClient("127.0.0.1:8080", grpc.WithStatsHandler(otelgrpc.NewClientHandler()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return message.NewMessageClient(conn)
})
