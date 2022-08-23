package client

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	grpci "github.com/actliboy/hoper/server/go/lib/utils/net/http/grpc"
	"github.com/actliboy/hoper/server/go/mod/protobuf/upload"
)

func GetUploadClient() upload.UploadServiceClient {
	// Set up a connection to the server.
	conn, err := grpci.GetDefaultClient("localhost:8090")

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	Connes = append(Connes, conn)
	return upload.NewUploadServiceClient(conn)
}
