package service

import (
	"context"
	"time"

	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/utils/time2"
)

type UserService struct {}

func (u *UserService) Signup(ctx context.Context, in *user.SignupReq) (*user.SignupRep, error) {
	return &user.SignupRep{Code:0,Data:&user.User{ActivatedAt:time2.Format(time.Now())},Msg:"test"}, nil
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
