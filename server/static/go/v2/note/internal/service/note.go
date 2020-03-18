package service

import (
	"context"
	"fmt"
	"unsafe"

	model "github.com/liov/hoper/go/v2/protobuf/note"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
)

type NoteService struct {
	model.UnimplementedNoteServiceServer
}

func (*NoteService) GC(ctx context.Context, req *model.GCReq) (*empty.Empty, error) {
	init := (*int)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	init1 := (*bool)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init1)
	init2 := (*byte)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init2)
	init3 := (*string)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init3)
	return &empty.Empty{}, nil
}
