package filter

import (
	"context"
	"strings"

	httpi "github.com/actliboy/hoper/server/go/lib/utils/net/http"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"google.golang.org/grpc"
)

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
