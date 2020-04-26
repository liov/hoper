package service

import "github.com/liov/hoper/go/v2/utils/net/http/pick"

func init() {
	pick.RegisterService(&UserService{}, &TestService{})
}
