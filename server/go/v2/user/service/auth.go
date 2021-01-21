package service

import (
	"context"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	"github.com/liov/hoper/go/v2/utils/log"
	jwti "github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"github.com/valyala/fasthttp"
)

func Auth(ctx *model.Ctx) (*model.AuthInfo, error) {
	if err := jwti.ParseToken(ctx.AuthInfo, ctx.Authorization, conf.Conf.Customize.TokenSecret); err != nil {
		return nil, err
	}
	conn := dao.NewUserRedis()
	defer conn.Close()
	err := conn.EfficientUserHashFromRedis(ctx.AuthInfo)
	if err != nil {
		log.Error(err)
		return nil, model.UserErr_InvalidToken
	}
	return ctx.AuthInfo, nil
}


type authKey struct{}

// AuthContext returns a new Context that carries value u.
func AuthContextF(r *fasthttp.Request) context.Context {
	return model.CtxFromContext(context.Background())
}

// FromContext returns the User value stored in ctx, if any.
func FromContextF(ctx context.Context) (*model.AuthInfo, bool) {
	user, ok := ctx.Value(authKey{}).(*model.AuthInfo)
	return user, ok
}