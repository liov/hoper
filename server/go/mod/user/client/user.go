package client

import (
	"crypto/tls"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	model "github.com/actliboy/hoper/server/go/mod/protobuf/user"
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
