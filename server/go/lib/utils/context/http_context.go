package contexti

import (
	"context"
	"encoding/base64"
	"github.com/google/uuid"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	"github.com/liov/hoper/server/go/lib/utils/net/http/request"
	timei "github.com/liov/hoper/server/go/lib/utils/time"
	"go.opencensus.io/trace"
	"go.opencensus.io/trace/propagation"
	"google.golang.org/grpc"
	"net/http"
	"net/url"
	"runtime"
	"strings"
	"sync"
	"time"
)

func GetPool[REQ any, P any]() sync.Pool {
	return sync.Pool{New: func() any {
		return new(RequestContext[REQ, P])
	}}
}

type AuthContext[REQ any, P any] interface {
	Authorization() error
	ParseRequest(r *http.Request) error
	RequestContext() *RequestContext[REQ, P]
}

type DeviceInfo struct {
	//设备
	Device     string `json:"device" gorm:"size:255"`
	Os         string `json:"os" gorm:"size:255"`
	AppCode    string `json:"appCode" gorm:"size:255"`
	AppVersion string `json:"appVersion" gorm:"size:255"`
	IP         string `json:"ip" gorm:"size:255"`
	Lng        string `json:"lng" gorm:"type:numeric(10,6)"`
	Lat        string `json:"lat" gorm:"type:numeric(10,6)"`
	Area       string `json:"area" gorm:"size:255"`
	UserAgent  string `json:"userAgent" gorm:"size:255"`
}

func DeviceFromHeader(r http.Header) *DeviceInfo {
	return Device(r.Get(httpi.HeaderDeviceInfo),
		r.Get(httpi.HeaderArea), r.Get(httpi.HeaderLocation),
		r.Get(httpi.HeaderUserAgent), r.Get(httpi.HeaderXForwardedFor))
}
func Device(infoHeader, area, localHeader, userAgent, ip string) *DeviceInfo {
	unknow := true
	var info DeviceInfo
	//Device-Info:device,osInfo,appCode,appVersion
	if infoHeader != "" {
		unknow = false
		var n, m int
		for i, c := range infoHeader {
			if c == '-' {
				switch n {
				case 0:
					info.Device = infoHeader[m:i]
				case 1:
					info.Os = infoHeader[m:i]
				case 2:
					info.AppCode = infoHeader[m:i]
				case 3:
					info.AppVersion = infoHeader[m:i]
				}
				m = i + 1
				n++
			}
		}
	}
	// area:xxx
	// location:1.23456,2.123456
	if area != "" {
		unknow = false
		info.Area, _ = url.PathUnescape(area)
	}
	if localHeader != "" {
		unknow = false
		var n, m int
		for i, c := range localHeader {
			if c == '-' {
				switch n {
				case 0:
					info.Lng = localHeader[m:i]
				case 1:
					info.Lat = localHeader[m:i]
				}
				m = i + 1
				n++
			}
		}

	}

	if userAgent != "" {
		unknow = false
		info.UserAgent = userAgent
	}
	if ip != "" {
		unknow = false
		info.IP = ip
	}
	if unknow {
		return nil
	}
	return &info
}

type RequestContext[REQ any, P any] struct {
	context.Context
	TraceID string
	Token   string
	*DeviceInfo
	request.RequestAt
	Request *REQ
	grpc.ServerTransportStream
	Internal string
	Values   map[string]any
	Props    P
}

func (c *RequestContext[REQ, P]) StartSpan(name string, o ...trace.StartOption) (*RequestContext[REQ, P], *trace.Span) {
	ctx, span := trace.StartSpan(c.Context, name, append(o, trace.WithSampler(trace.AlwaysSample()),
		trace.WithSpanKind(trace.SpanKindServer))...)
	c.Context = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID.String()
	}
	return c, span
}

func CtxWithRequest[P any](r *http.Request, tracing bool) (*RequestContext[http.Request, P], *trace.Span) {
	var span *trace.Span
	ctx := context.Background()
	if r != nil {
		ctx = r.Context()
		if tracing {
			// go.opencensus.io/trace 完全包含了golang.org/x/net/trace 的功能
			// grpc内置配合,看了源码并没有启用，根本没调用
			// 系统trace只能追踪单个请求，且只记录时间及是否完成，只能/debug/requests看
			/*			t = gtrace.New(methodFamily(r.RequestURI), r.RequestURI)
						ctx = gtrace.NewContext(ctx, t)
			*/

			// 直接从远程读取Trace信息，Trace是否为空交给propagation包判断
			traceString := r.Header.Get(httpi.GrpcTraceBin)
			if traceString == "" {
				traceString = r.Header.Get(httpi.HeaderTraceBin)
			}
			var traceBin []byte
			if len(traceString)%4 == 0 {
				// Input was padded, or padding was not necessary.
				traceBin, _ = base64.StdEncoding.DecodeString(traceString)
			}
			traceBin, _ = base64.RawStdEncoding.DecodeString(traceString)

			if parent, ok := propagation.FromBinary(traceBin); ok {
				ctx, span = trace.StartSpanWithRemoteParent(ctx, r.RequestURI,
					parent, trace.WithSampler(trace.AlwaysSample()),
					trace.WithSpanKind(trace.SpanKindServer))
			} else {
				ctx, span = trace.StartSpan(ctx, r.RequestURI,
					trace.WithSampler(trace.AlwaysSample()),
					trace.WithSpanKind(trace.SpanKindServer))
			}
		}
	} else {
		if tracing {
			pc, _, _, _ := runtime.Caller(2)
			f := runtime.FuncForPC(pc)

			/*			t = gtrace.New(file, fmt.Sprintf("%s:%d", file, line))
						ctx = gtrace.NewContext(ctx, t)*/

			ctx, span = trace.StartSpan(ctx, f.Name(),
				trace.WithSampler(trace.AlwaysSample()),
				trace.WithSpanKind(trace.SpanKindServer))
		}
	}

	ctxi := NewCtx[http.Request, P](ctx)
	setWithHttpReq(ctxi, r)
	return ctxi, span
}

func methodFamily(m string) string {
	m = strings.TrimPrefix(m, "/") // remove leading slash
	if i := strings.Index(m, "/"); i >= 0 {
		m = m[:i] // remove everything from second slash
	}
	return m
}

type ctxKey struct{}

func (c *RequestContext[REQ, P]) ContextWrapper() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, c)
}

func CtxFromContext[REQ any, P any](ctx context.Context) *RequestContext[REQ, P] {
	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*RequestContext[REQ, P])
	if !ok {
		c = NewCtx[REQ, P](ctx)
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	return c
}

func (c *RequestContext[REQ, P]) WithContext(ctx context.Context) {
	c.Context = ctx
}

func NewCtx[REQ any, P any](ctx context.Context) *RequestContext[REQ, P] {
	span := trace.FromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID.String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	return &RequestContext[REQ, P]{
		Context: ctx,
		TraceID: traceId,
		RequestAt: request.RequestAt{
			Time:       now,
			TimeStamp:  now.Unix(),
			TimeString: now.Format(timei.TimeFormat),
		},
		ServerTransportStream: grpc.ServerTransportStreamFromContext(ctx),
	}
}

func setWithHttpReq[P any](c *RequestContext[http.Request, P], r *http.Request) {
	if r == nil {
		return
	}
	c.Request = r
	c.DeviceInfo = DeviceFromHeader(r.Header)
	c.Internal = r.Header.Get(httpi.GrpcInternal)
	c.Token = httpi.GetToken(r)
}

func (c *RequestContext[REQ, P]) reset(ctx context.Context) *RequestContext[REQ, P] {
	span := trace.FromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID.String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	c.Context = ctx
	c.RequestAt.Time = now
	c.RequestAt.TimeString = now.Format(timei.TimeFormat)
	c.RequestAt.TimeStamp = now.Unix()
	return c
}

type HttpContext[P any] RequestContext[http.Request, P]
