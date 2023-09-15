package service

import (
	"context"
	"google.golang.org/protobuf/types/known/wrapperspb"

	model "github.com/liovx/hoper/server/go/protobuf/content"
)

type NoteService struct {
	model.UnimplementedNoteServiceServer
}

func (*NoteService) Create(ctx context.Context, req *model.Note) (*wrapperspb.StringValue, error) {
	return &wrapperspb.StringValue{Value: "成功"}, nil
}
