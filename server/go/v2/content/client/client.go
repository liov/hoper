package client

import (
	user_model "github.com/liov/hoper/go/v2/protobuf/user"
	"google.golang.org/grpc"
)

var (
	Connes     clientConns
	UserClient user_model.UserServiceClient
)

func init() {
	var conn *grpc.ClientConn
	UserClient, conn = GetUserClient()
	Connes = append(Connes, conn)
}


type clientConns []*grpc.ClientConn

func (cs clientConns) Close() {
	for _, conn := range cs {
		conn.Close()
	}
}
