package service

import (
	"context"
	"time"

	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/dao"
	"github.com/liov/hoper/go/v2/utils/time2"
	"github.com/liov/hoper/go/v2/utils/valid"
)

type UserService struct{}

func (*UserService) Signup(ctx context.Context, in *user.SignupReq) (*user.SignupRep, error) {
	u := user.User{ActivatedAt: time2.Format(time.Now()), Birthday: time2.Format(time.Now())}
	err := valid.Validate.Struct(&u)
	if err != nil {
		return &user.SignupRep{Code:1,Data:nil,Msg:valid.Trans(err)}, nil
	}
	dao.Dao.DB.Create(&u)
	return &user.SignupRep{Code: 0, Data: &u, Msg: "test"}, nil
}

func (*UserService) Edit(context.Context, *user.EditReq) (*user.EditRep, error) {
	panic("implement me")
}

func (*UserService) Login(context.Context, *user.LoginReq) (*user.LoginRep, error) {
	panic("implement me")
}

func (*UserService) Logout(context.Context, *user.LogoutReq) (*user.LogoutRep, error) {
	panic("implement me")
}

func (*UserService) GetUser(context.Context, *user.GetReq) (*user.User, error) {
	panic("implement me")
}
