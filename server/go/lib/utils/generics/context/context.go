package contexti

import (
	"context"
	contexti "github.com/actliboy/hoper/server/go/lib/utils/context"
	"github.com/actliboy/hoper/server/go/lib/utils/net/http/request"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc"
	"net/http"
)

// TODO
type RequestContext[REQ http.Request | fasthttp.Request, P any] struct {
	context.Context
	TraceID string
	Token   string
	*contexti.DeviceInfo
	request.RequestAt
	Request *REQ
	grpc.ServerTransportStream
	Internal string
	Values   map[string]interface{}
	Props    P
}
