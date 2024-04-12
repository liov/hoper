package service

import (
	"github.com/hopeio/cherry/context/http_context"
	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/service"
)

func auth(ctx *http_context.Context, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
