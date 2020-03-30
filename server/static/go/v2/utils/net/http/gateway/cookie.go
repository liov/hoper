package gateway

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

/*
* grpc的本地代理实现相当的有问题
* 凡是local_request_中的metadata都没有赋值操作，
* 然后还有个ctx = runtime.NewServerMetadataContext(ctx, md)
* 然后在ForwardResponseMessage中有各种关于对空metadata的操作
* 根本没有作用的调用，这是极大的性能浪费，也没法借此设置返回的的header
 */

func CookieHook(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	if v, ok := message.(SetCookie); ok {
		writer.Header().Add("Set-Cookie", v.GetCookie())
	}
	return nil
}

func GrpcSetCookie(ctx context.Context, cookie string) {
	err := grpc.SendHeader(ctx, metadata.MD{"Set-Cookie": []string{cookie}})
	if err != nil {
		md, _ := metadata.FromIncomingContext(ctx)
		md.Set("Set-Cookie", cookie)
	}
}

type SetCookie interface {
	GetCookie() string
}
