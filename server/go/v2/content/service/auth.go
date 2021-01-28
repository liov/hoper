package service

import (
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
)

func Auth(ctx *model.Ctx) (*model.AuthInfo, error) {
	if err := ctx.ParseToken(ctx.Authorization, conf.Conf.Customize.TokenSecret); err != nil {
		return nil, err
	}
	conn := dao.NewUserRedis()
	defer conn.Close()
	err := conn.EfficientUserHashFromRedis(ctx)
	if err != nil {
		return nil, model.UserErr_InvalidToken
	}
	return ctx.AuthInfo, nil
}


