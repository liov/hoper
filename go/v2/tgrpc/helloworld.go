package helloworld

import (
	"context"
	"net"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/initialize/dao"
	pb "github.com/liov/hoper/go/v2/protobuf"
	"github.com/liov/hoper/go/v2/user/internal/config"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func main() {
	defer log.Sync()
	defer dao.Dao.Close()
	initialize.Start(config.Conf)
	log.Info(*dao.Dao)
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
