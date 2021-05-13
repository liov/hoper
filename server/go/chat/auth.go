package chat

import (
	"github.com/liov/hoper/v2/protobuf/user"
	contexti "github.com/liov/hoper/v2/tiga/context"
	"github.com/liov/hoper/v2/user/service"
)

func auth(ctx *contexti.Ctx, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
