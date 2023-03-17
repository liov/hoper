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

type Ctx = contexti.Ctx[fasthttp.Request, uint8]

func CtxWithRequest(ctx context.Context, r *fasthttp.Request) *Ctx {
	ctxi := contexti2.NewCtx[fasthttp.Request, uint8](ctx)
	c := &Ctx{RequestContext: ctxi}
	setWithReq(c, r)
	return c
}

func setWithReq(c *Ctx, r *fasthttp.Request) {
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
