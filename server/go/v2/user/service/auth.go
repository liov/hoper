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


func Auth(ctx *model.Ctx) (*model.AuthInfo, error) {
	parts := strings.Split(ctx.Authorization, ".")
	if len(parts) != 3 {
		return nil, model.UserErr_InvalidToken
	}
	signature := parts[2]
	cacheTmp, err := dao.Dao.Cache.Get(signature)
	if err == nil {
		if cache, ok := cacheTmp.(*model.Cache); ok {
			cache.LastActiveAt = ctx.RequestUnix
			ctx.AuthInfo = cache.AuthInfo
			ctx.Authorization = cache.Authorization
			return cache.AuthInfo, nil
		}
	}
	if err := ctx.ParseToken(ctx.Authorization, conf.Conf.Customize.TokenSecret); err != nil {
		return nil, err
	}
	err = userRedis.EfficientUserHashFromRedis(ctx)
	if err != nil {
		return nil, model.UserErr_InvalidToken
	}
	dao.Dao.Cache.SetWithExpire(signature,
		&model.Cache{AuthInfo: ctx.AuthInfo, Authorization: ctx.Authorization},
		5*time.Second)
	return ctx.AuthInfo, nil
}

// AuthContext returns a new Context that carries value u.
func FasthttpCtx(r *fasthttp.Request) pick.Context {
	ctx := model.CtxFromContext(context.Background())
	ctx.Authorization = fasthttpi.GetToken(r)
	return ctx
}
