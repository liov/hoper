package service

import (
	"context"
	"github.com/liov/hoper/go/v2/protobuf/user"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	"strings"
	"time"

	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	fasthttpi "github.com/liov/hoper/go/v2/utils/net/http/fasthttp"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"

	"github.com/valyala/fasthttp"
)

var ExportAuth = auth

func auth(ctx *contexti.Ctx, update bool) (*user.AuthInfo, error) {
	signature := ctx.Token[strings.LastIndexByte(ctx.Token, '.')+1:]
	cacheTmp, ok := dao.Dao.Cache.Get(signature)
	if ok {
		cache := cacheTmp.(*contexti.Authorization)
		cache.LastActiveAt = ctx.TimeStamp
		ctx.Authorization = cache
		auth := cache.AuthInfo.(*user.AuthInfo)
		return auth, nil
	}
	auth := &user.AuthInfo{}
	ctx.AuthInfo = auth
	if err := ctx.ParseToken(ctx.Token, conf.Conf.Customize.TokenSecret); err != nil {
		return nil, err
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

// AuthContext returns a new Context that carries value u.
func FasthttpCtx(r *fasthttp.Request) pick.Context {
	ctx := contexti.CtxFromContext(context.Background())
	ctx.Token = fasthttpi.GetToken(r)
	return ctx
}
