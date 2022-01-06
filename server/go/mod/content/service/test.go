package service

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/actliboy/hoper/server/go/lib/protobuf/empty"
	"github.com/actliboy/hoper/server/go/lib/tiga"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	model "github.com/actliboy/hoper/server/go/mod/protobuf/content"
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
