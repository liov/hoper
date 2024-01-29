package rpc

import (
	"github.com/hopeio/lemon/utils/log"
	grpci "github.com/hopeio/lemon/utils/net/http/grpc"
	"github.com/liov/hoper/server/go/protobuf/upload"
	"github.com/liov/hoper/server/go/protobuf/user"
)

var (
	UserClient   user.UserServiceClient
	UploadClient upload.UploadServiceClient
)

func init() {
	UserClient = GetUserClient()
	UploadClient = GetUploadClient()
}

func GetUserClient() user.UserServiceClient {
	// Set up a connection to the server.
	conn, err := grpci.GetDefaultClient("localhost:8090")
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return user.NewUserServiceClient(conn)
}

func GetUploadClient() upload.UploadServiceClient {
	// Set up a connection to the server.
	conn, err := grpci.GetDefaultClient("localhost:8090")

	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	return upload.NewUploadServiceClient(conn)
}
