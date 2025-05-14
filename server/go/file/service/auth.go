package service

import (
	"github.com/hopeio/context/httpctx"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/service"
)

func auth(ctx *httpctx.Context, update bool) (*user.AuthBase, error) {
	return service.ExportAuth(ctx, update)
}
