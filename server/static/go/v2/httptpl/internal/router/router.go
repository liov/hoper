package router

import (
	"github.com/liov/hoper/go/v2/httptpl/internal/controller"
	"github.com/liov/hoper/go/v2/utils/http/iris_plus"
)

func route() []iris_plus.Controller {
	//controller注册
	return []iris_plus.Controller{
		&controller.UserController{},
	}
}
