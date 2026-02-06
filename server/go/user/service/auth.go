package service

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	jwt2 "github.com/hopeio/scaffold/jwt"
	"github.com/liov/hoper/server/go/protobuf/user"

	"github.com/liov/hoper/server/go/user/data"
	"github.com/liov/hoper/server/go/user/global"
)

var ExportAuth = auth
var jwtValidator = jwt.NewValidator()

func auth(ctx context.Context, update bool) (*user.AuthInfo, error) {
	token := ctx.Auth().Token
	signature := token[strings.LastIndexByte(token, '.')+1:]
	cacheTmp, ok := global.Dao.Cache.Get(signature)
	if ok {
		cache := cacheTmp.(*jwt2.Claims[*user.AuthInfo])
		err := jwtValidator.Validate(cache)
		if err != nil {
			return nil, err
		}
		return cache.Auth, nil
	}
	authorization, err := jwt2.Auth[*user.AuthInfo](ctx, global.Conf.User.TokenSecretBytes)
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
