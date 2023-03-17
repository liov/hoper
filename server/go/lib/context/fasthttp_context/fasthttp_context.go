package fasthttp_context

import (
	"context"
	contexti "github.com/liov/hoper/server/go/lib/context"
	contexti2 "github.com/liov/hoper/server/go/lib/utils/context"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	fasthttpi "github.com/liov/hoper/server/go/lib/utils/net/http/fasthttp"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/valyala/fasthttp"
)

type Context = contexti.Context[fasthttp.Request]

func ContextWithRequest(ctx context.Context, r *fasthttp.Request) *Context {
	ctxi := contexti2.NewCtx[fasthttp.Request](ctx)
	c := &Context{RequestContext: ctxi}
	setWithReq(c, r)
	return c
}

func setWithReq(c *Context, r *fasthttp.Request) {
	c.Request = r
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = Device(&r.Header)
	c.Internal = stringsi.ToString(r.Header.Peek(httpi.GrpcInternal))
}

func Device(r *fasthttp.RequestHeader) *contexti2.DeviceInfo {
	return contexti2.Device(stringsi.ToString(r.Peek(httpi.HeaderDeviceInfo)),
		stringsi.ToString(r.Peek(httpi.HeaderArea)),
		stringsi.ToString(r.Peek(httpi.HeaderLocation)),
		stringsi.ToString(r.Peek(httpi.HeaderUserAgent)),
		stringsi.ToString(r.Peek(httpi.HeaderXForwardedFor)),
	)
}
