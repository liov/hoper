package client

import (
	"github.com/actliboy/hoper/server/go/mod/protobuf/upload"
	"github.com/actliboy/hoper/server/go/mod/protobuf/user"
	"google.golang.org/grpc"
)

var (
	Connes       clientConns
	UserClient   user.UserServiceClient
	UploadClient upload.UploadServiceClient
)

func init() {
	UserClient = GetUserClient()
	UploadClient = GetUploadClient()
}

type clientConns []*grpc.ClientConn

func (cs clientConns) Close() {
	for _, conn := range cs {
		conn.Close()
	}
}
