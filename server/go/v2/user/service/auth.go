package service

import (
	"context"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	fasthttpi "github.com/liov/hoper/go/v2/utils/net/fasthttp"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	jwti "github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"github.com/valyala/fasthttp"
)

func init() {
	jwt.WithUnmarshaller(model.JWTUnmarshaller)(jwti.Parser)
}

func auth(ctx *model.Ctx, update bool) error {
	signature := ctx.Token[strings.LastIndexByte(ctx.Token, '.')+1:]
	cacheTmp, err := dao.Dao.Cache.Get(signature)
	if err == nil {
		if cache, ok := cacheTmp.(*model.Authorization); ok {
			cache.LastActiveAt = ctx.RequestUnix
			ctx.Authorization = cache
			return nil
		}
	}
	if err := ctx.ParseToken(ctx.Token, conf.Conf.Customize.TokenSecret); err != nil {
		return err
	}
	ctx.LastActiveAt = ctx.RequestUnix
	if update {
		err = userRedis.EfficientUserHashFromRedis(ctx)
		if err != nil {
			return model.UserErrInvalidToken
		}
	}
	dao.Dao.Cache.SetWithExpire(signature, ctx.Authorization, 5*time.Second)
	return nil
}

func Auth(ctx *model.Ctx) error {
	return auth(ctx, false)
}

func AuthWithUpdate(ctx *model.Ctx) error {
	return auth(ctx, true)
}

// AuthContext returns a new Context that carries value u.
func FasthttpCtx(r *fasthttp.Request) pick.Context {
	ctx := model.CtxFromContext(context.Background())
	ctx.Token = fasthttpi.GetToken(r)
	return ctx
}
