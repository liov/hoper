package service

import (
	"context"

	model "github.com/actliboy/hoper/server/go/protobuf/content"
	"github.com/hopeio/pandora/protobuf/response"
)

type NoteService struct {
	model.UnimplementedNoteServiceServer
}

func (*NoteService) Create(ctx context.Context, req *model.Note) (*response.CommonRep, error) {
	return &response.CommonRep{Message: "成功"}, nil
}
