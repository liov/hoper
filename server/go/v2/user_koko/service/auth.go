package service

import (
	"context"

	model "github.com/liov/hoper/go/v2/protobuf/user"
	fasthttpi "github.com/liov/hoper/go/v2/utils/net/fasthttp"
	"github.com/liov/hoper/go/v2/utils/net/http/koko"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	"github.com/valyala/fasthttp"
)

func Auth(ctx *koko.Ctx) error {
	return userRedis.EfficientUserHashFromRedis(ctx)
}

// AuthContext returns a new Context that carries value u.
func FasthttpCtx(r *fasthttp.Request) pick.Context {
	ctx := model.CtxFromContext(context.Background())
	ctx.Authorization = fasthttpi.GetToken(r)
	return ctx
}
