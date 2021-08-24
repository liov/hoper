package grpci

import (
	"github.com/liov/hoper/server/go/lib/utils/log"
	"github.com/liov/hoper/server/go/lib/utils/net/http/grpc/stats"
	"google.golang.org/grpc"
)

func GetDefaultClient(target string) *grpc.ClientConn {
	// Set up a connection to the server.
	conn, err := grpc.Dial(target, grpc.WithInsecure(),
		grpc.WithStatsHandler(&stats.ClientHandler{}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn
}
