package data

import (
	"github.com/liov/hoper/server/go/protobuf/user"
	"google.golang.org/protobuf/types/known/emptypb"

	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

func TestVerifyCode(t *testing.T) {
	md := runtime.ServerMetadata{}
	u, err := UserClient().VerifyCode(metadata.NewOutgoingContext(context.Background(),
		metadata.MD{"key": []string{"value"}}),
		&emptypb.Empty{},
		grpc.Header(&md.HeaderMD),
		grpc.Trailer(&md.TrailerMD),
	)
	if err != nil {
		t.Log(err)
	}
	t.Log(u)
}

func TestSendVerifyCode(t *testing.T) {
	md := runtime.ServerMetadata{}
	_, err := UserClient().SendVerifyCode(metadata.NewOutgoingContext(context.Background(),
		metadata.MD{"key": []string{"value"}}),
		&user.SendVerifyCodeReq{},
		grpc.Header(&md.HeaderMD),
		grpc.Trailer(&md.TrailerMD),
	)
	if err != nil {
		t.Log(err)
	}
}
