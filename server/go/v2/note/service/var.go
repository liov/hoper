package service

import (
	"github.com/liov/hoper/go/v2/note/internal/dao"
)

var (
	NoteSvc = &NoteService{}

	userDao = &dao.NoteDao{}
)
