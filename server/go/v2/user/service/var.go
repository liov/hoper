package service

import (
	"github.com/liov/hoper/go/v2/user/dao"
)

var (
	userSvc  = &UserService{}
	oauthSvc *OauthService

	userDao = &dao.UserDao{}
)

func GetUserService() *UserService {
	if userSvc != nil {
		return userSvc
	}
	userSvc = new(UserService)
	return userSvc
}