package contexti

import (
	"context"
	fasthttpi "github.com/liov/hoper/v2/utils/net/http/fasthttp"
	"github.com/liov/hoper/v2/utils/net/http/request"
	stringsi "github.com/liov/hoper/v2/utils/strings"
	timei "github.com/liov/hoper/v2/utils/time"
	"github.com/valyala/fasthttp"
	"sync"
	"time"

	"github.com/google/uuid"

	httpi "github.com/liov/hoper/v2/utils/net/http"

	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)

var (
	fastCtxPool = sync.Pool{New: func() interface{} {
		return new(FastCtx)
	}}
)

type FastCtx struct {
	context.Context
	TraceID string
	*Authorization
	*DeviceInfo
	request.RequestAt
	Request *fasthttp.Request
	grpc.ServerTransportStream
	Internal string
}

func (c *FastCtx) StartSpan(name string, o ...trace.StartOption) (*FastCtx, *trace.Span) {
	ctx, span := trace.StartSpan(c.Context, name, o...)
	c.Context = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID.String()
	}
	return c, span
}

func (c *FastCtx) WithContext(ctx context.Context) {
	c.Context = ctx
}

func (ctxi *FastCtx) ContextWrapper() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, ctxi)
}

func CtxWithFastRequest(ctx context.Context, r *fasthttp.Request) *FastCtx {
	ctxi := newFastCtx(ctx)
	ctxi.setWithReq(r)
	return ctxi
}

func CtxFromFastRequest(r *fasthttp.Request) *FastCtx {
	ctxi := newFastCtx(context.Background())
	ctxi.setWithReq(r)
	return ctxi
}

func FastCtxFromContext(ctx context.Context) *FastCtx {
	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*FastCtx)
	if !ok {
		c = newFastCtx(ctx)
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	return c
}

func newFastCtx(ctx context.Context) *FastCtx {
	span := trace.FromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID.String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	return &FastCtx{
		Context:       ctx,
		TraceID:       traceId,
		Authorization: &Authorization{},
		RequestAt: request.RequestAt{
			Time:       now,
			TimeStamp:  now.Unix(),
			TimeString: now.Format(timei.FormatTime),
		},
	}
}

func (c *FastCtx) setWithReq(r *fasthttp.Request) {
	c.Request = r
	c.Token = fasthttpi.GetToken(r)
	c.DeviceInfo = DeviceFast(&r.Header)
	c.Internal = stringsi.ToString(r.Header.Peek(httpi.GrpcInternal))
}

func (c *FastCtx) reset(ctx context.Context) *FastCtx {
	span := trace.FromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID.String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	c.Context = ctx
	c.RequestAt.Time = now
	c.RequestAt.TimeString = now.Format(timei.FormatTime)
	c.RequestAt.TimeStamp = now.Unix()
	return c
}

func (c *FastCtx) GetAuthInfo(auth func(*FastCtx) error) (AuthInfo, error) {
	if c.Authorization == nil {
		c.Authorization = new(Authorization)
	}
	if err := auth(c); err != nil {
		return nil, err
	}
	return c.AuthInfo, nil
}

func DeviceFast(r *fasthttp.RequestHeader) *DeviceInfo {
	return device(stringsi.ToString(r.Peek(httpi.HeaderDeviceInfo)),
		stringsi.ToString(r.Peek(httpi.HeaderArea)),
		stringsi.ToString(r.Peek(httpi.HeaderLocation)),
		stringsi.ToString(r.Peek(httpi.HeaderUserAgent)),
		stringsi.ToString(r.Peek(httpi.HeaderXForwardedFor)),
	)
}
