package client

import (
	"github.com/liov/hoper/go/v2/protobuf/user"
	"google.golang.org/grpc"
)

func init() {
	client,conn  := GetUserClient()
	Client.UserClient = client
	Client.conns = append(Client.conns,conn)
}

var Client client
type client struct {
	UserClient user.UserServiceClient
	conns []*grpc.ClientConn
}

func (c *client) Close() {
	for _,conn:=range c.conns{
		conn.Close()
	}
}