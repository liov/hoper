package upload

import (
	contexti "github.com/actliboy/hoper/server/go/lib/tiga/context"
	"github.com/actliboy/hoper/server/go/mod/protobuf/user"
	"github.com/actliboy/hoper/server/go/mod/user/service"
)

func auth(ctx *contexti.Ctx, update bool) (*user.AuthInfo, error) {
	return service.ExportAuth(ctx, update)
}
