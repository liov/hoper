package gateway

import (
	"context"
	"net/http"
	"net/url"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/encoding/protobuf/jsonpb"
	"github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/request"
	"google.golang.org/grpc/metadata"
)

type GatewayHandle func(context.Context, *runtime.ServeMux)

func Gateway(gatewayHandle GatewayHandle) http.Handler {
	ctx := context.Background()

	gwmux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &jsonpb.JSONPb{API: json.Standard}),

		runtime.WithMetadata(func(ctx context.Context, req *http.Request) metadata.MD {
			area, err := url.PathUnescape(req.Header.Get(request.Area))
			if err != nil {
				area = ""
			}
			var token = httpi.GetToken(req)

			return map[string][]string{
				request.Area:          {area},
				request.DeviceInfo:    {req.Header.Get(request.DeviceInfo)},
				request.Location:      {req.Header.Get(request.Location)},
				request.Authorization: {token},
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
	runtime.WithRoutingErrorHandler(RoutingErrorHandler)(gwmux)
	runtime.WithErrorHandler(CustomHTTPError)(gwmux)
	if gatewayHandle != nil {
		gatewayHandle(ctx, gwmux)
	}
	return gwmux
}
