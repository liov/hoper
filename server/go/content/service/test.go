package service

import (
	"context"
	"fmt"
	"google.golang.org/protobuf/types/known/emptypb"
	"unsafe"

	model "github.com/actliboy/hoper/server/go/protobuf/content"
)

type TestService struct {
	model.UnimplementedTestServiceServer
}

func (*TestService) GC(ctx context.Context, req *model.GCReq) (*emptypb.Empty, error) {
	//address:= strconv.FormatUint()
	init := (*TestService)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	return &emptypb.Empty{}, nil
}

func (*TestService) Restart(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}
