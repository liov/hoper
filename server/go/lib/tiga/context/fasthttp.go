package contexti

import (
	"context"
	contexti "github.com/liov/hoper/server/go/lib/utils/context"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	fasthttpi "github.com/liov/hoper/server/go/lib/utils/net/http/fasthttp"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/valyala/fasthttp"
)

func CtxWithFastRequest(ctx context.Context, r *fasthttp.Request) *Ctx {
	ctxi := contexti.NewCtx(ctx)
	c := &Ctx{RequestContext: ctxi}
	c.setWithFastReq(r)
	return c
}

func (c *Ctx) setWithFastReq(r *fasthttp.Request) {
	c.FastRequest = r
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = DeviceFast(&r.Header)
	c.Internal = stringsi.ToString(r.Header.Peek(httpi.GrpcInternal))
}

func DeviceFast(r *fasthttp.RequestHeader) *contexti.DeviceInfo {
	return contexti.Device(stringsi.ToString(r.Peek(httpi.HeaderDeviceInfo)),
		stringsi.ToString(r.Peek(httpi.HeaderArea)),
		stringsi.ToString(r.Peek(httpi.HeaderLocation)),
		stringsi.ToString(r.Peek(httpi.HeaderUserAgent)),
		stringsi.ToString(r.Peek(httpi.HeaderXForwardedFor)),
	)
}
