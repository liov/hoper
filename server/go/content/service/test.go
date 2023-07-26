package service

import (
	"context"
	"fmt"
	"unsafe"

	model "github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/hopeio/zeta/protobuf/empty"
)

type TestService struct {
	model.UnimplementedTestServiceServer
}

func (*TestService) GC(ctx context.Context, req *model.GCReq) (*empty.Empty, error) {
	//address:= strconv.FormatUint()
	init := (*TestService)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	return &empty.Empty{}, nil
}

func (*TestService) Restart(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	return &empty.Empty{}, nil
}
