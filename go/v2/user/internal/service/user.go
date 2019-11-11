package service

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"time"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/internal/config"
	modelconst "github.com/liov/hoper/go/v2/user/internal/model"
	"github.com/liov/hoper/go/v2/utils/json/protobuf"
	"github.com/liov/hoper/go/v2/utils/strings2"
	"github.com/liov/hoper/go/v2/utils/time2"
	"github.com/liov/hoper/go/v2/utils/valid"
)

type UserService struct{}

func (*UserService) Signup(ctx context.Context, in *model.SignupReq) (*model.SignupRep, error) {
	var rep = &model.SignupRep{Code: 10000}
	err := valid.Validate.Struct(in)
	if err != nil {
		rep.Msg = valid.Trans(err)
		return rep, nil
	}

	if in.Email == "" && in.Phone == "" {
		rep.Msg = "请填写邮箱或手机号"
		return rep, err
	}

	if exist, err := userDao.ExitByEmailORPhone(in.Email, in.Phone); exist {
		if  in.Email != "" {
			rep.Msg = "邮箱已被注册"
			return rep, err
		} else  {
			rep.Msg = "手机号已被注册"
			return rep, err
		}
	}
	var user = &model.User{}
	user.Name = in.Name
	user.Email = in.Email
	user.Gender = modelconst.UserSexNil
	user.CreatedAt = time2.Format(time.Now())
	user.UpdatedAt = user.CreatedAt
	user.Status = modelconst.UserStatusInActive
	user.AvatarURL = modelconst.DefaultAvatar
	user.Password = encryptPassword(in.Password)
	if err = userDao.Creat(user); err != nil {
		rep.Msg = "新建出错"
		return rep, err
	}
	data,_:=json.Marshal(user)
	rep.Data = &protobuf.StringJson{S:string(data)}
	rep.Msg = "新建成功"
	return rep, nil
}

// Salt 每个用户都有一个不同的盐
func salt(password string) string {
	return password[0:5]
}

// EncryptPassword 给密码加密
func encryptPassword(password string) string {
	hash := salt(password) + config.Conf.Server.PassSalt + password[5:]
	return fmt.Sprintf("%x",md5.Sum(strings2.ToBytes(hash)))
}

//验证密码是否正确
func checkPassword(password string, user *model.User) bool {
	if password == "" || user.Password == "" {
		return false
	}
	return encryptPassword(password) == user.Password
}

func (*UserService) Edit(context.Context, *model.EditReq) (*model.EditRep, error) {
	panic("implement me")
}

func (*UserService) Login(context.Context, *model.LoginReq) (*model.LoginRep, error) {
	panic("implement me")
}

func (*UserService) Logout(context.Context, *model.LogoutReq) (*model.LogoutRep, error) {
	panic("implement me")
}

func (*UserService) GetUser(context.Context, *model.GetReq) (*model.User, error) {
	panic("implement me")
}
