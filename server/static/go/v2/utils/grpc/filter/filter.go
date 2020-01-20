package filter

import (
	"context"
	"runtime/debug"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"

	"github.com/liov/hoper/go/v2/utils/log"
	"google.golang.org/grpc"
)

func filter(
	ctx context.Context, req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.CallTwo.Errorf("%v panic: %v", info, r)
			debug.PrintStack()
			err = errorcode.SysError.Err()
		}
		//不能添加错误处理，除非所有返回的结构相同
		if err != nil {
			if errcode, ok := err.(errorcode.ErrCode); ok {
				err = errcode.Err()
			}
		}
	}()

	return handler(ctx, req)
}

func CommonUnaryServerInterceptor() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{
		//filter应该在最前
		filter, grpc_validator.UnaryServerInterceptor(),
	}
}

func CommonStreamServerInterceptor() []grpc.StreamServerInterceptor {
	return []grpc.StreamServerInterceptor{
		grpc_validator.StreamServerInterceptor(),
	}
}
