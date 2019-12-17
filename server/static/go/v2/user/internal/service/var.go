package service

import "github.com/liov/hoper/go/v2/user/internal/dao"

var (
	UserSvc = &UserService{}

	userDao = &dao.UserDao{}
)
