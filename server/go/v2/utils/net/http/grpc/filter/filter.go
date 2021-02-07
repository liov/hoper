package filter

import (
	"context"
	"fmt"
	"strings"

	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/validator"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	runtimei "github.com/liov/hoper/go/v2/utils/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

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
			frame,_:=runtimei.GetCallerFrame(2)
			log.Default.With(zap.String(log.Stack, fmt.Sprintf("%s:%d (%#x)\n\t%s\n", frame.File, frame.Line, frame.PC, frame.Function))).Error(" panic: ", r)
			err = errorcode.SysError.ErrRep()
		}
		//不能添加错误处理，除非所有返回的结构相同
		if err != nil {
			if _, ok := err.(interface{GRPCStatus() *status.Status}); !ok {
				err = errorcode.Unknown.Message(err.Error())
			}
		}
	}()

	return handler(ctx, req)
}

func UnaryServerInterceptor(i ...grpc.UnaryServerInterceptor) []grpc.UnaryServerInterceptor {
	return append([]grpc.UnaryServerInterceptor{
		//filter应该在最前
		filter, grpc_validator.UnaryServerInterceptor(),
	}, i...)
}

func StreamServerInterceptor(i ...grpc.StreamServerInterceptor) []grpc.StreamServerInterceptor {
	return append([]grpc.StreamServerInterceptor{
		grpc_validator.StreamServerInterceptor(),
	}, i...)
}

//鉴权过程大体和自定义的一致，就不用中间件的形式了
func ensureValidToken(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.InvalidArgument, "missing metadata")
	}
	// The keys within metadata.Header are normalized to lowercase.
	// See: https://godoc.org/google.golang.org/grpc/metadata#New
	if !Valid(md[httpi.HeaderAuthorization]) {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token")
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}

//可自定义
var Valid = func(authorization []string) bool {
	if len(authorization) < 1 {
		return false
	}
	token := strings.TrimPrefix(authorization[0], "Bearer ")
	// Perform the token validation here. For the sake of this example, the code
	// here forgoes any of the usual OAuth2 token validation and instead checks
	// for a token matching an arbitrary string.
	if token != "some-secret-token" {
		return false
	}
	return true
}

type authority struct {
	FunctionalAuthority bool
	DataAuthority       interface{}
}
