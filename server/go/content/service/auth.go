package service

import (
	"context"

	"github.com/liov/hoper/server/go/user/model"
	"github.com/liov/hoper/server/go/user/service"
)

func auth(ctx context.Context, update bool) (*model.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
