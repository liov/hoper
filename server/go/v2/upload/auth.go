package upload

import (
	"github.com/liov/hoper/go/v2/protobuf/user"
	contexti "github.com/liov/hoper/go/v2/tailmon/context"
	"github.com/liov/hoper/go/v2/user/service"
)

func auth(ctx *contexti.Ctx, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx,update)
}