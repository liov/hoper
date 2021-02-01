package service

import (
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/service"
)

func Auth(ctx *model.Ctx) error {
	return service.Auth(ctx)
}

func AuthWithUpdate(ctx *model.Ctx) error {
	return service.AuthWithUpdate(ctx)
}

