package gateway

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"
	ihttp "github.com/liov/hoper/go/v2/protobuf/utils/http"
)

func ResponseHook(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	if res, ok := message.(*ihttp.Response); ok {
		for k, v := range res.Header {
			writer.Header().Add(k, v)
		}
	}
	return nil
}
