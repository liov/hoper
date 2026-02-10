package service

import (
	"context"

	"github.com/liov/hoper/server/go/protobuf/user"
	"github.com/liov/hoper/server/go/user/service"
)

func auth(ctx context.Context, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
