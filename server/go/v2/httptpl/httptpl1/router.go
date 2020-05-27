package main

import (
	"github.com/liov/hoper/go/v2/httptpl/httptpl1/internal/controller"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"
)

func init() {
	api.RegisterService(&controller.UserController{}, "user")
}
