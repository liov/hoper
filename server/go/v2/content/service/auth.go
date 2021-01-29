package service

import (
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/service"
)

func Auth(ctx *model.Ctx) (*model.AuthInfo, error) {
	return service.Auth(ctx)
}


