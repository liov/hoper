package client

import (
	"github.com/liov/hoper/v2/protobuf/utils/empty"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestGetUserClient(t *testing.T) {
	md := runtime.ServerMetadata{}
	user, err := UserClient.VerifyCode(metadata.NewOutgoingContext(context.Background(),
		metadata.MD{"key": []string{"value"}}),
		&empty.Empty{},
		grpc.Header(&md.HeaderMD),
		grpc.Trailer(&md.TrailerMD),
	)
	if err != nil {
		t.Log(err)
	}
	t.Log(user)
}
