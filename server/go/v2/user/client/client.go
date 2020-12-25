package client

import (
	user_model "github.com/liov/hoper/go/v2/protobuf/user"
	"google.golang.org/grpc"
)

var (
	UserClient user_model.UserServiceClient
)

func init() {
	var conn *grpc.ClientConn
	UserClient, conn = GetUserClient()
	ClientConns = append(ClientConns, conn)
}

var ClientConns clientConns

type clientConns []*grpc.ClientConn

func (cs clientConns) Close() {
	for _, conn := range cs {
		conn.Close()
	}
}
