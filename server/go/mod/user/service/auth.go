package service

import (
	contexti "github.com/liov/hoper/server/go/lib/context"
	"github.com/liov/hoper/server/go/lib/context/http_context"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
	"strings"
	"time"

	"github.com/liov/hoper/server/go/mod/user/conf"
	"github.com/liov/hoper/server/go/mod/user/dao"
)

var ExportAuth = auth

func auth(ctx *http_context.Context, update bool) (*user.AuthInfo, error) {
	signature := ctx.Token[strings.LastIndexByte(ctx.Token, '.')+1:]
	cacheTmp, ok := dao.Dao.Cache.Get(signature)
	if ok {
		cache := cacheTmp.(*contexti.Authorization)
		ctx.Props.LastActiveAt = ctx.TimeStamp
		ctx.Props.Authorization = cache
		auth := cache.AuthInfo.(*user.AuthInfo)
		return auth, nil
	}
	auth := &user.AuthInfo{}
	ctx.Props.AuthInfo = auth
	if err := ctx.Props.ParseToken(ctx.Token, conf.Conf.Customize.TokenSecret); err != nil {
		return nil, user.UserErrLoginTimeout
	}
	ctx.Props.LastActiveAt = ctx.TimeStamp
	if update {
		userDao := dao.GetDao(ctx)
		err := userDao.EfficientUserHashFromRedis()
		if err != nil {
			return nil, err
		}
	}
	if !ok {
		dao.Dao.Cache.SetWithTTL(signature, ctx.Props.Authorization, 0, 5*time.Second)
	}
	return auth, nil
}
