package client

import (
	"crypto/tls"
	"github.com/hopeio/pandora/utils/log"
	model "github.com/liov/hoper/server/go/protobuf/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func GetUserClient() (model.UserServiceClient, *grpc.ClientConn) {
	// Set up a connection to the server.
	conn, err := grpc.Dial("grpc.hoper.xyz:443", grpc.WithTransportCredentials(credentials.NewTLS(&tls.Config{ServerName: "grpc.hoper.xyz", InsecureSkipVerify: true})))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return model.NewUserServiceClient(conn), conn
}
