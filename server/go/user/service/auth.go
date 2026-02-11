package service

import (
	"context"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hopeio/cherry"
	httpx "github.com/hopeio/gox/net/http"
	contextx "github.com/hopeio/scaffold/context"
	jwt2 "github.com/hopeio/scaffold/jwt"
	"github.com/liov/hoper/server/go/protobuf/user"

	"github.com/liov/hoper/server/go/user/data"
	"github.com/liov/hoper/server/go/user/global"
)

var ExportAuth = auth
var jwtValidator = jwt.NewValidator()

func auth(ctx context.Context, update bool) (*user.AuthInfo, error) {
	metadata := cherry.GetMetadata(ctx)
	if metadata == nil {
		return nil, user.UserErrNoLogin
	}
	metadata.Token = httpx.GetToken(metadata.Request.Header)
	token := metadata.Token
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
		rclient := global.Dao.Redis.Client.WithContext(ctx)
		userDao := data.GetRedisDao(rclient)
		err := userDao.EfficientUserHashFromRedis(ctx, authorization.Auth)
		if err != nil {
			return nil, err
		}
	}
	if metadata.Data == nil {
		metadata.Data = &user.ClientInfo{
			Auth: authorization.Auth,
		}
	} else {
		metadata.Data.(*user.ClientInfo).Auth = authorization.Auth
	}

	global.Dao.Cache.SetWithTTL(signature, authorization, 0, 5*time.Second)
	return authorization.Auth, nil
}

func Device(ctx context.Context) *contextx.DeviceInfo {
	metadata := cherry.GetMetadata(ctx)
	var device *contextx.DeviceInfo
	if metadata.Data != nil {
		clientInfo := metadata.Data.(*user.ClientInfo)
		if clientInfo.Device != nil {
			return clientInfo.Device
		} else {
			device = contextx.DeviceFromHeader(metadata.Request.Header)
			clientInfo.Device = device
		}
	} else {
		device = contextx.DeviceFromHeader(metadata.Request.Header)
		metadata.Data = &user.ClientInfo{
			Device: device,
		}
	}
	return device
}
