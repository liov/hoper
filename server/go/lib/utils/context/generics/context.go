package generics

import (
	"context"
	contexti "github.com/actliboy/hoper/server/go/lib/utils/context"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/request"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"net/http"
)

// TODO
type RequestContext[T http.Request | fasthttp.Request] struct {
	context.Context
	TraceID string
	Token   string
	*contexti.DeviceInfo
	request.RequestAt
	Request *T
	grpc.ServerTransportStream
	Internal string
	Values   map[string]interface{}
}
