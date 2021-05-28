package contexti

import (
	"context"
	httpi "github.com/liov/hoper/v2/utils/net/http"
	fasthttpi "github.com/liov/hoper/v2/utils/net/http/fasthttp"
	stringsi "github.com/liov/hoper/v2/utils/strings"
	"github.com/valyala/fasthttp"
)

func CtxWithFastRequest(ctx context.Context, r *fasthttp.Request) *Ctx {
	ctxi := newCtx(ctx)
	ctxi.setWithFastReq(r)
	return ctxi
}

func CtxFromFastRequest(r *fasthttp.Request) *Ctx {
	ctxi := newCtx(context.Background())
	ctxi.setWithFastReq(r)
	return ctxi
}

func (c *Ctx) setWithFastReq(r *fasthttp.Request) {
	c.FastRequest = r
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = DeviceFast(&r.Header)
	c.Internal = stringsi.ToString(r.Header.Peek(httpi.GrpcInternal))
}

func DeviceFast(r *fasthttp.RequestHeader) *DeviceInfo {
	return device(stringsi.ToString(r.Peek(httpi.HeaderDeviceInfo)),
		stringsi.ToString(r.Peek(httpi.HeaderArea)),
		stringsi.ToString(r.Peek(httpi.HeaderLocation)),
		stringsi.ToString(r.Peek(httpi.HeaderUserAgent)),
		stringsi.ToString(r.Peek(httpi.HeaderXForwardedFor)),
	)
}
