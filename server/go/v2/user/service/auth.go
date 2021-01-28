package service

import (
	"context"

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

// AuthContext returns a new Context that carries value u.
func FasthttpCtx(r *fasthttp.Request) pick.Context {
	ctx:=model.CtxFromContext(context.Background())
	ctx.Authorization = fasthttpi.GetToken(r)
	return ctx
}