package service

import (
	"github.com/gorilla/sessions"
	model "github.com/liov/hoper/go/v2/protobuf/user"
)

type UserService struct{}

func (*UserService) Add(ctx *sessions.Session, req *model.SignupReq) (*model.SignupRep, error) {
	return &model.SignupRep{Message: "测试"}, nil
}

func (*UserService) Edit(ctx *sessions.Session, req *model.EditReq) (*model.EditReq_EditDetails, error) {
	return nil, nil
}
