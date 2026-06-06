package service

import (
	"context"

	"github.com/liov/hoper/server/go/user/service"
	"github.com/liov/hoper/server/go/user/model"
)

func auth(ctx context.Context, update bool) (*model.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
