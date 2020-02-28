package gateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/liov/hoper/go/v2/utils/json"
	"github.com/liov/hoper/go/v2/utils/protobuf/jsonpb"
	"google.golang.org/grpc/metadata"
)

func Gateway(gatewayHandle func(context.Context, *runtime.ServeMux)) http.Handler {
	ctx := context.Background()

	jsonpb := &jsonpb.JSONPb{
		json.Json,
	}

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, jsonpb),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
		runtime.WithMetadata(func(ctx context.Context, request *http.Request) metadata.MD {
			area, err := url.PathUnescape(request.Header.Get("area"))
			if err != nil {
				area = ""
			}
			var auth string
			cookie, _ := request.Cookie("token")
			if cookie != nil {
				value, _ := url.QueryUnescape(cookie.Value)
				if value == "" {
					auth = request.Header.Get("authorization")
				} else {
					auth = value
				}
			}

			return map[string][]string{
				"device-info": {request.Header.Get("device-info")},
				"location":    {area, request.Header.Get("location")},
				"auth":        {auth},
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
	if gatewayHandle != nil {
		gatewayHandle(ctx, gwmux)
	}
	return gwmux
}
