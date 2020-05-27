package filter

import (
	"context"
	"github.com/liov/hoper/go/v2/utils/verification/luosimao"
	"google.golang.org/grpc/metadata"
)

func LuosimaoVerify(reqURL, apiKey string, ctx context.Context) error {
	md, _ := metadata.FromIncomingContext(ctx)
	response := md.Get("luosimao")
	if len(response) == 0 || response[0] == "" {
		return luosimao.LuosimaoErr
	}
	return luosimao.LuosimaoVerify(reqURL, apiKey, response[0])
}
