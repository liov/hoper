package service

import (
	contexti "github.com/actliboy/hoper/server/go/lib/context"
	"github.com/actliboy/hoper/server/go/mod/protobuf/user"
	"strings"
	"time"

	"github.com/actliboy/hoper/server/go/mod/user/conf"
	"github.com/actliboy/hoper/server/go/mod/user/dao"
)

var ExportAuth = auth

func auth(ctx *contexti.Ctx, update bool) (*user.AuthInfo, error) {
	signature := ctx.Token[strings.LastIndexByte(ctx.Token, '.')+1:]
	cacheTmp, ok := dao.Dao.Cache.Get(signature)
	if ok {
		cache := cacheTmp.(*contexti.Authorization)
		ctx.LastActiveAt = ctx.TimeStamp
		ctx.Authorization = cache
		auth := cache.AuthInfo.(*user.AuthInfo)
		return auth, nil
	}
	auth := &user.AuthInfo{}
	ctx.AuthInfo = auth
	if err := ctx.ParseToken(ctx.Token, conf.Conf.Customize.TokenSecret); err != nil {
		return nil, user.UserErrLoginTimeout
	}
	ctx.LastActiveAt = ctx.TimeStamp
	if update {
		userDao := dao.GetDao(ctx)
		err := userDao.EfficientUserHashFromRedis()
		if err != nil {
			return nil, err
		}
	}
	if !ok {
		dao.Dao.Cache.SetWithTTL(signature, ctx.Authorization, 0, 5*time.Second)
	}
	return auth, nil
}
