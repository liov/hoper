package router

import (
	"github.com/liov/hoper/go/v2/httptpl/internal/controller"
	"github.com/liov/hoper/go/v2/httptpl/internal/service"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"
)

func Route() []api.Controller {
	//controller注册
	return []api.Controller{
		&controller.UserController{Service: &service.UserService{}},
	}
}
