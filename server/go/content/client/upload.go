package client

import (
	"github.com/liov/hoper/v2/protobuf/upload"
	"github.com/liov/hoper/v2/utils/log"
	"github.com/liov/hoper/v2/utils/net/http/grpc/stats"
	"google.golang.org/grpc"
)

func GetUploadClient() upload.UploadServiceClient {
	// Set up a connection to the server.
	conn, err := grpc.Dial("localhost:8090", grpc.WithInsecure(),
		grpc.WithStatsHandler(&stats.ClientHandler{}))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	Connes = append(Connes, conn)
	return upload.NewUploadServiceClient(conn)
}
