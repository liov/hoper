package main

import (
	"context"
	"os"
	"time"

	"github.com/liov/hoper/go/v2/gateway/internal/config"
	pb "github.com/liov/hoper/go/v2/gateway/protobuf"
	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/initialize/dao"
	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc"
)

const (
	address     = "localhost:9090"
	defaultName = "world"
)

//验证，无论是在一个大模块里还是拆分成小模块，都是可以rpc通信的
func main() {
	defer log.Sync()
	defer dao.Dao.Close()
	initialize.Start(config.Conf)
	log.Info(*dao.Dao)
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewSayClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Hello(ctx, &pb.Request{Name: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Infof("Greeting: %s", r.Msg)
}
