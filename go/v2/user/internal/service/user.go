package service

import (
	"context"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/liov/hoper/go/v2/protobuf/user"
)

type UserService struct {}

func (u *UserService) Signup(ctx context.Context, in *user.SignupReq) (*user.SignupRep, error) {
	now:=time.Now()
	return &user.SignupRep{Code:0,Data:&user.User{ActivatedAt:&timestamp.Timestamp{Seconds:int64(now.Second()),Nanos:int32(now.Nanosecond())}},Msg:"test"}, nil
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
