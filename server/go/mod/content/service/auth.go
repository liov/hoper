package service

import (
	"github.com/liov/hoper/server/go/lib/context/http_context"
	"github.com/liov/hoper/server/go/mod/protobuf/user"
	"github.com/liov/hoper/server/go/mod/user/service"
)

func auth(ctx *http_context.Context, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
