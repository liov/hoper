package service

import (
	"github.com/liov/hoper/go/v2/user_koko/dao"
)

var (
	userSvc  = &UserService{}
	oauthSvc *OauthService

	userDao = &dao.UserDao{}
	userRedis = &dao.UserRedis{}
)

func GetUserService() *UserService {
	if userSvc != nil {
		return userSvc
	}
	userSvc = new(UserService)
	return userSvc
}