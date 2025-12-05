package service

import (
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hopeio/gox/context/httpctx"
	jwt2 "github.com/hopeio/scaffold/jwt"
	"github.com/liov/hoper/server/go/protobuf/user"

	"github.com/liov/hoper/server/go/user/data"
	"github.com/liov/hoper/server/go/user/global"
)

var ExportAuth = auth
var jwtValidator = jwt.NewValidator()

func auth(ctx *httpctx.Context, update bool) (*user.AuthBase, error) {
	signature := ctx.Token[strings.LastIndexByte(ctx.Token, '.')+1:]
	cacheTmp, ok := global.Dao.Cache.Get(signature)
	if ok {
		cache := cacheTmp.(*jwt2.Claims[*user.AuthBase])
		err := jwtValidator.Validate(cache)
		if err != nil {
			return nil, err
		}
		return cache.Auth, nil
	}
	authorization, err := jwt2.Auth[httpctx.RequestCtx, *user.AuthBase](ctx, global.Conf.User.TokenSecretBytes)
	if err != nil {
		return nil, user.UserErrNoLogin
	}

	if update {
		userDao := data.GetRedisDao(ctx, global.Dao.Redis.Client)
		err := userDao.EfficientUserHashFromRedis()
		if err != nil {
			return nil, err
		}
	}
	global.Dao.Cache.SetWithTTL(signature, authorization, 0, 5*time.Second)
	return authorization.Auth, nil
}
