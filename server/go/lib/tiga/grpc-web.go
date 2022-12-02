package tiga

import (
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"net/http"
)

func NewGrpcWebServer(grpcServer *grpc.Server) *grpcweb.WrappedGrpcServer {
	return grpcweb.WrapServer(grpcServer, grpcweb.WithAllowedRequestHeaders([]string{"*"}), grpcweb.WithWebsockets(true), grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
		return true
	}), grpcweb.WithOriginFunc(func(origin string) bool {
		return true
	}))
}
