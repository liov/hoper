package oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/oauth"
	"github.com/liov/hoper/go/v2/utils/encoding/schema"
	"github.com/liov/hoper/go/v2/utils/net/http"
	"google.golang.org/grpc/metadata"
)

func RegisterOauthServiceHandlerServer(r *gin.Engine, server user.OauthServiceServer) {
	r.GET("/oauth/authorize", func(ctx *gin.Context) {
		token := httpi.GetToken(ctx.Request)
		var protoReq oauth.OauthReq
		schema.DefaultDecoder.Decode(&protoReq, ctx.Request.URL.Query())
		res, _ := server.OauthAuthorize(
			metadata.NewIncomingContext(
				ctx.Request.Context(),
				metadata.MD{"auth": {token}}),
			&protoReq)

		res.Response(ctx.Writer)
	})

	r.POST("/oauth/access_token", func(ctx *gin.Context) {
		var protoReq oauth.OauthReq
		schema.DefaultDecoder.Decode(&protoReq, ctx.Request.PostForm)
		res, _ := server.OauthToken(ctx.Request.Context(), &protoReq)
		res.Response(ctx.Writer)
	})
}
