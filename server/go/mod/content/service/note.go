package service

import (
	"context"

	"github.com/liov/hoper/server/go/lib/protobuf/response"
	model "github.com/liov/hoper/server/go/mod/protobuf/content"
)

type NoteService struct {
	model.UnimplementedNoteServiceServer
}

func (*NoteService) Create(ctx context.Context, req *model.Note) (*response.CommonRep, error) {
	return &response.CommonRep{Message: "成功"}, nil
}
