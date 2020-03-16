package service

import model "github.com/liov/hoper/go/v2/protobuf/note"

type NoteService struct {
	model.UnimplementedNoteServiceServer
}
