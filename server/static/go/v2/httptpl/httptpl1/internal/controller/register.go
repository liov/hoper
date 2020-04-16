package controller

import "github.com/liov/hoper/go/v2/utils/net/http/iris/api"

func init() {
	api.RegisterService(&UserController{}, "user")
}
