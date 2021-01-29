package oauth

import (
	"github.com/gin-gonic/gin"
	"github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/protobuf/utils/oauth"
	"github.com/liov/hoper/go/v2/utils/net/http/request/binding"

	"github.com/liov/hoper/go/v2/utils/net/http"
	"google.golang.org/grpc/metadata"
)

func RegisterOauthServiceHandlerServer(r *gin.Engine, server user.OauthServiceServer) {
	r.GET("/oauth/authorize", func(ctx *gin.Context) {
		var protoReq oauth.OauthReq
		binding.DefaultDecoder().Decode(&protoReq, ctx.Request.URL.Query())
		res, _ := server.OauthAuthorize(
			metadata.NewIncomingContext(
				ctx.Request.Context(),
				metadata.MD{"auth": {httpi.GetToken(ctx.Request)}}),
			&protoReq)

		res.Response(ctx.Writer)
	})

	r.POST("/oauth/access_token", func(ctx *gin.Context) {
		var protoReq oauth.OauthReq
		binding.DefaultDecoder().Decode(&protoReq, ctx.Request.PostForm)
		res, _ := server.OauthToken(ctx.Request.Context(), &protoReq)
		res.Response(ctx.Writer)
	})
}
