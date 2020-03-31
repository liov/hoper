package router

import (
	"github.com/liov/hoper/go/v2/httptpl2/internal/service"
	"github.com/liov/hoper/go/v2/utils/net/http/iris/api"
)

func Route() []api.Service {
	//controller注册
	return []api.Service{
		&service.UserService{},
	}
}
