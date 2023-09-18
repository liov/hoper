package service

import (
	contexti "github.com/hopeio/lemon/context"
	"github.com/hopeio/lemon/context/http_context"
	"github.com/liovx/hoper/server/go/protobuf/user"
	"strings"
	"time"

	"github.com/liovx/hoper/server/go/user/confdao"
	"github.com/liovx/hoper/server/go/user/data"
)

var ExportAuth = auth

func auth(ctx *http_context.Context, update bool) (*user.AuthInfo, error) {
	signature := ctx.Token[strings.LastIndexByte(ctx.Token, '.')+1:]
	cacheTmp, ok := confdao.Dao.Cache.Get(signature)
	if ok {
		cache := cacheTmp.(*contexti.Authorization)
		ctx.LastActiveAt = ctx.TimeStamp
		ctx.Authorization = cache
		auth := cache.AuthInfo.(*user.AuthInfo)
		return auth, nil
	}
	auth := &user.AuthInfo{}
	ctx.AuthInfo = auth
	if err := ctx.ParseToken(ctx.Token, confdao.Conf.Customize.TokenSecret); err != nil {
		return nil, user.UserErrLoginTimeout
	}
	ctx.LastActiveAt = ctx.TimeStamp
	if update {
		userDao := data.GetDao(ctx)
		err := userDao.EfficientUserHashFromRedis()
		if err != nil {
			return nil, err
		}
	}
	if !ok {
		confdao.Dao.Cache.SetWithTTL(signature, ctx.Authorization, 0, 5*time.Second)
	}
	return auth, nil
}
