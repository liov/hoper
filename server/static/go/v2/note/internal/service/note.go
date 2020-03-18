package service

import (
	"context"
	"fmt"
	"unsafe"

	v2 "github.com/liov/hoper/go/v2/initialize/v2"
	model "github.com/liov/hoper/go/v2/protobuf/note"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
)

type NoteService struct {
	model.UnimplementedNoteServiceServer
}

func (*NoteService) GC(ctx context.Context, req *model.GCReq) (*empty.Empty, error) {
	init := (*v2.InitConfig)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	return &empty.Empty{}, nil
}
