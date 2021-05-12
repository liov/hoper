package service

import (
	"context"
	"fmt"
	"unsafe"

	model "github.com/liov/hoper/go/v2/protobuf/content"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/tiga"
	"github.com/liov/hoper/go/v2/tiga/initialize"
)

type TestService struct {
	model.UnimplementedTestServiceServer
}

func (*TestService) GC(ctx context.Context, req *model.GCReq) (*empty.Empty, error) {
	//address:= strconv.FormatUint()
	init := (*initialize.Init)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	return &empty.Empty{}, nil
}

func (*TestService) Restart(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	tiga.ReStart()
	return &empty.Empty{}, nil
}
