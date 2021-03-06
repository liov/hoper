package client

import (
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/protobuf/utils/request"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGetUserClient(t *testing.T) {
	md:=runtime.ServerMetadata{}
	user,err:=UserClient.AuthInfo(metadata.NewOutgoingContext(context.Background(), metadata.MD{"key":[]string{"value"}}),
		&empty.Empty{},
		grpc.Header(&md.HeaderMD),
		grpc.Trailer(&md.TrailerMD),
		)
	if err!=nil{
		t.Log(err)
	}
	t.Log(user)
}
