package service

import (
	"context"

	"github.com/liov/hoper/go/v2/protobuf/user"
)

type UserService struct {}

func (u *UserService) Signup(ctx context.Context, in *user.SignupReq) (*user.LoginRep, error) {
	return &user.LoginRep{}, nil
}

func (u *UserService) Edit(context.Context, *user.EditReq) (*user.EditRep, error) {
	panic("implement me")
}

func (u *UserService) Login(context.Context, *user.LoginReq) (*user.LoginRep, error) {
	panic("implement me")
}

func (u *UserService) Logout(context.Context, *user.LogoutReq) (*user.LogoutRep, error) {
	panic("implement me")
}

func (u *UserService) GetUser(context.Context, *user.GetReq) (*user.User, error) {
	panic("implement me")
}
