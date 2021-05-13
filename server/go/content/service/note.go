package service

import (
	"context"

	model "github.com/liov/hoper/v2/protobuf/content"
	"github.com/liov/hoper/v2/protobuf/utils/response"
)

type NoteService struct {
	model.UnimplementedNoteServiceServer
}

func (*NoteService) Create(ctx context.Context, req *model.Note) (*response.CommonRep, error) {
	return &response.CommonRep{Message: "成功"}, nil
}
