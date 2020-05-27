package oauth

import (
	"github.com/kataras/iris/v12"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/oauth"
	"github.com/liov/hoper/go/v2/utils/encoding/schema"
	"github.com/liov/hoper/go/v2/utils/net/http/auth"
	"google.golang.org/grpc/metadata"
)

func RegisterOauthServiceHandlerServer(mux *iris.Application, server user.OauthServiceServer) {
	mux.Get("/oauth/authorize", func(ctx iris.Context) {
		token := auth.GetToken(ctx.Request())
		var protoReq oauth.OauthReq
		schema.DefaultDecoder.Decode(&protoReq, ctx.Request().URL.Query())
		res, _ := server.OauthAuthorize(
			metadata.NewIncomingContext(
				ctx.Request().Context(),
				metadata.MD{"auth": {token}}),
			&protoReq)

		res.Response(ctx.ResponseWriter())
	})

	mux.Post("/oauth/access_token", func(ctx iris.Context) {
		var protoReq oauth.OauthReq
		schema.DefaultDecoder.Decode(&protoReq, ctx.FormValues())
		res, _ := server.OauthToken(ctx.Request().Context(), &protoReq)
		res.Response(ctx.ResponseWriter())
	})
}
