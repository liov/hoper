package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/liov/hoper/go/v2/protobuf"
	"github.com/liov/hoper/go/v2/utils/time2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

//go:generate protoc -I../../../../../proto/ -I$GOPATH/src -I$GOPATH/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis -I$GOPATH/src/github.com/gogo/protobuf/protobuf  ../../../../../proto/*.proto --gogo_out=plugins=grpc,Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types,Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api,Mgoogle/protobuf/field_mask.proto=github.com/gogo/protobuf/types,Mgoogle/protobuf/any.proto=github.com/gogo/protobuf/types:../protobuf

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterGreeterServer(s, &Service{})
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

type Service struct{}

func (*Service) SayHello(ctx context.Context, in *protobuf.HelloRequest) (*protobuf.HelloReply, error) {
	t := time.Now()
	tt := time2.Time2(t)
	log.Println("请求： ", in, int64(tt), t.Unix())
	return &protobuf.HelloReply{Time: tt}, nil
}
