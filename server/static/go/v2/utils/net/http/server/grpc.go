package server

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/filter"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func (s *Server) Grpc() *grpc.Server {
	if s.GRPCRegistr == nil {
		return nil
	}
	gs := grpc.NewServer(
		//filter应该在最前
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(filter.CommonUnaryServerInterceptor()...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(filter.CommonStreamServerInterceptor()...)),
	)
	s.GRPCRegistr(gs)
	// Register reflection service on gRPC server.
	reflection.Register(gs)
	return gs
}
