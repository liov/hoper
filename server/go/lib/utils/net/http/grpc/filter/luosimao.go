package filter

import (
	"context"
	"github.com/actliboy/hoper/server/go/lib/utils/verification"
	"google.golang.org/grpc/metadata"
)

func LuosimaoVerify(reqURL, apiKey string, ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	response := md.Get("luosimao")
	if len(response) == 0 || response[0] == "" {
		return verification.LuosimaoErr
	}
	return verification.LuosimaoVerify(reqURL, apiKey, response[0])
}
