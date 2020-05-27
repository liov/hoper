package service

import (
	"github.com/liov/hoper/go/v2/user/internal/dao"
)

var (
	userSvc  = &UserService{}
	oauthSvc *OauthService

	userDao = &dao.UserDao{}
)
