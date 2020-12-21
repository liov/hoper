package gateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/encoding/protobuf/jsonpb"
	"github.com/liov/hoper/go/v2/utils/net/http/auth"
	"google.golang.org/grpc/metadata"
)

type GatewayHandle func(context.Context, *runtime.ServeMux)

func Gateway(gatewayHandle GatewayHandle) http.Handler {
	ctx := context.Background()

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &jsonpb.JSONPb{API: json.Standard}),
		//runtime.WithProtoErrorHandler(CustomHTTPError),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			area, err := url.PathUnescape(request.Header.Get("area"))
			if err != nil {
				area = ""
			}
			var token = auth.GetToken(request)

			return map[string][]string{
				"device-info": {request.Header.Get("device-info")},
				"location":    {area, request.Header.Get("location")},
				"auth":        {token},
			}
		}),
		runtime.WithIncomingHeaderMatcher(func(key string) (string, bool) {
			switch key {
			case
				"Accept",
				"Accept-Charset",
				"Accept-Language",
				"Accept-Ranges",
				//"Authorization",
				"Cache-Control",
				"Content-Type",
				//"Cookie",
				"Date",
				"Expect",
				"From",
				"Host",
				"If-Match",
				"If-Modified-Since",
				"If-None-Match",
				"If-Schedule-Tag-Match",
				"If-Unmodified-Since",
				"Max-Forwards",
				"Origin",
				"Pragma",
				"Referer",
				"User-Agent",
				"Via",
				"Warning":
				return key, true
			}
			return "", false
		}),
		runtime.WithOutgoingHeaderMatcher(
			func(key string) (string, bool) {
				switch key {
				case
					"set-cookie":
					return key, true
				}
				return "", false
			}))

	runtime.WithForwardResponseOption(CookieHook)(gwmux)
	runtime.WithForwardResponseOption(ResponseHook)(gwmux)
	if gatewayHandle != nil {
		gatewayHandle(ctx, gwmux)
	}
	return gwmux
}
