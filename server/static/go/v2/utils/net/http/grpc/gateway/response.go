package gateway

import (
	"context"
	"net/http"

	"github.com/golang/protobuf/proto"
	"github.com/liov/hoper/go/v2/protobuf/utils/response"
)

func ResponseHook(ctx context.Context, writer http.ResponseWriter, message proto.Message) error {
	if res, ok := message.(*response.HttpResponse); ok {
		for k, v := range res.Header {
			writer.Header().Add(k, v)
		}
		writer.WriteHeader(int(res.StatusCode))
	}
	return nil
}
