package main

import (
	"context"
	"os"
	"time"

	pb "github.com/liov/hoper/go/v2/protobuf"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc"
)

const (
	address     = "172.27.168.110:50051"
	defaultName = "world"
)

//验证，无论是在一个大模块里还是拆分成小模块，都是可以rpc通信的
func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Infof("Greeting: %s", r.Message)
}
