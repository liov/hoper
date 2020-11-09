package tailmon

import (
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/liov/hoper/go/v2/utils/net/http/grpc/filter"
	"google.golang.org/grpc"
)

/*如果grpc给了ServerOptions单独的设置api这个方法就有用了
func (s *Server) Grpc() *grpc.Server {
	if s.GRPCRegistr == nil {
		return nil
	}
	gs := grpc.NewServer(
		//filter应该在最前
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(filter.UnaryServerInterceptor()...)),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(filter.StreamServerInterceptor()...)),
	)
	s.GRPCRegistr(gs)
	// Register reflection service on gRPC server.
	reflection.Register(gs)
	initialize.BasicConfig.Register()
	return gs
}
*/

func DefaultGRPCServer() *grpc.Server {
	gs := grpc.NewServer(
		//filter应该在最前
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				filter.UnaryServerInterceptor()...,
			)),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				filter.StreamServerInterceptor()...,
			)),
	)
	return gs
}
