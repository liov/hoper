package service

import (
	"github.com/hopeio/lemon/context/http_context"
	"github.com/liovx/hoper/server/go/protobuf/user"
	"github.com/liovx/hoper/server/go/user/service"
)

func auth(ctx *http_context.Context, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
