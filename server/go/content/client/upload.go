package client

import (
	"github.com/actliboy/hoper/server/go/protobuf/upload"
	"github.com/hopeio/pandora/utils/log"
	grpci "github.com/hopeio/pandora/utils/net/http/grpc"
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
