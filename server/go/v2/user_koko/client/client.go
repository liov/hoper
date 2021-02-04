package client

import (
	user_model "github.com/liov/hoper/go/v2/protobuf/user"
	"google.golang.org/grpc"
)

var (
	Connes     clientConnes
	UserClient user_model.UserServiceClient
)

func init() {
	var conn *grpc.ClientConn
	UserClient, conn = GetUserClient()
	Connes = append(Connes, conn)
}

type clientConnes []*grpc.ClientConn

func (cs clientConnes) Close() {
	for _, conn := range cs {
		conn.Close()
	}
}
