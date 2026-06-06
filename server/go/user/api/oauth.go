/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package api

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/gox/strstruct"
	"github.com/hopeio/protobuf/response"
	"github.com/liov/hoper/server/go/protobuf/user"

	httpx "github.com/hopeio/gox/net/http"

	"google.golang.org/grpc/metadata"
)

type OauthServiceServer interface {
	OauthAuthorize(context.Context, *user.OauthReq) (*response.HttpResponse, error)
	OauthToken(context.Context, *user.OauthReq) (*response.HttpResponse, error)
}

func RegisterOauthServiceHandlerServer(r *gin.Engine, server OauthServiceServer) {
	r.GET("/oauth/authorize", func(ctx *gin.Context) {
		var protoReq user.OauthReq
		strstruct.DefaultDecoder().Decode(&protoReq, ctx.Request.URL.Query())
		res, _ := server.OauthAuthorize(
			metadata.NewIncomingContext(
				ctx.Request.Context(),
				metadata.MD{"auth": {httpx.GetToken(ctx.Request.Header)}}),
			&protoReq)

		res.Respond(ctx, ctx.Writer)
	})

	r.POST("/oauth/access_token", func(ctx *gin.Context) {
		var protoReq user.OauthReq
		strstruct.DefaultDecoder().Decode(&protoReq, ctx.Request.PostForm)
		res, _ := server.OauthToken(ctx.Request.Context(), &protoReq)
		res.Respond(ctx, ctx.Writer)
	})
}
