package service

import (
	"context"
	"fmt"
	"unsafe"

	"github.com/liov/hoper/go/v2/initialize/v2"
	model "github.com/liov/hoper/go/v2/protobuf/note"
	"github.com/liov/hoper/go/v2/protobuf/utils/empty"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
	"github.com/liov/hoper/go/v2/utils/server"
)

type NoteService struct {
	model.UnimplementedNoteServiceServer
}

func (*NoteService) Create(ctx context.Context, req *model.Note) (*response.CommonRep, error) {
	return &response.CommonRep{Message: "成功"}, nil
}

func (*NoteService) GC(ctx context.Context, req *model.GCReq) (*empty.Empty, error) {
	//address:= strconv.FormatUint()
	init := (*initialize.InitConfig)(unsafe.Pointer(uintptr(req.Address)))
	fmt.Println(*init)
	return &empty.Empty{}, nil
}

func (*NoteService) Restart(ctx context.Context, req *empty.Empty) (*empty.Empty, error) {
	server.ReStart()
	return &empty.Empty{}, nil
}
