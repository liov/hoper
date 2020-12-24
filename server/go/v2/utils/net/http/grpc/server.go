package grpc

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/filter"
	"google.golang.org/grpc"
)


func DefaultGRPCServer(usi []grpc.UnaryServerInterceptor,ssi []grpc.StreamServerInterceptor) *grpc.Server {
	gs := grpc.NewServer(
		//filter应该在最前
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				filter.UnaryServerInterceptor(usi...)...,
			)),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				filter.StreamServerInterceptor(ssi...)...,
			)),
	)
	return gs
}

