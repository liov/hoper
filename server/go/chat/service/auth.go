package service

import (
	"github.com/actliboy/hoper/server/go/protobuf/user"
	"github.com/actliboy/hoper/server/go/user/service"
	"github.com/hopeio/zeta/context/http_context"
)

func auth(ctx *http_context.Context, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
