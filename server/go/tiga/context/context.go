package contexti

import (
	"context"
	"errors"
	"github.com/liov/hoper/v2/utils/net/http/request"
	timei "github.com/liov/hoper/v2/utils/time"
	"github.com/valyala/fasthttp"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/liov/hoper/v2/utils/encoding/json"
	httpi "github.com/liov/hoper/v2/utils/net/http"

	stringsi "github.com/liov/hoper/v2/utils/strings"
	jwti "github.com/liov/hoper/v2/utils/verification/auth/jwt"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
)

var (
	ctxPool = sync.Pool{New: func() interface{} {
		return new(Ctx)
	}}
)

type AuthInfo interface {
	IdStr() string
}

type Authorization struct {
	AuthInfo     `json:"auth"`
	IdStr        string `json:"-" gorm:"-"`
	LastActiveAt int64  `json:"lat,omitempty"`
	ExpiredAt    int64  `json:"exp,omitempty"`
	LoginAt      int64  `json:"iat,omitempty"`
	Token        string `json:"-"`
}

func (x *Authorization) Valid(helper *jwt.ValidationHelper) error {
	if x.ExpiredAt != 0 && x.LastActiveAt > x.ExpiredAt {
		return errors.New("登录过期")
	}
	return nil
}

func (x *Authorization) GenerateToken(secret []byte) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func (x *Authorization) ParseToken(token, secret string) error {
	if err := jwti.ParseToken(x, token, secret); err != nil {
		return err
	}
	x.IdStr = x.AuthInfo.IdStr()
	return nil
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

func Device(r http.Header) *DeviceInfo {
	return device(r.Get(httpi.HeaderDeviceInfo),
		r.Get(httpi.HeaderArea), r.Get(httpi.HeaderLocation),
		r.Get(httpi.HeaderUserAgent), r.Get(httpi.HeaderXForwardedFor))
}
func device(infoHeader, area, localHeader, userAgent, ip string) *DeviceInfo {
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

type Ctx struct {
	context.Context
	TraceID string
	*Authorization
	*DeviceInfo
	request.RequestAt
	Request     *http.Request
	FastRequest *fasthttp.Request
	grpc.ServerTransportStream
	Internal string
}

func (c *Ctx) StartSpan(name string, o ...trace.StartOption) (*Ctx, *trace.Span) {
	ctx, span := trace.StartSpan(c.Context, name, o...)
	c.Context = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID.String()
	}
	return c, span
}

func (c *Ctx) WithContext(ctx context.Context) {
	c.Context = ctx
}

type ctxKey struct{}

func (ctxi *Ctx) ContextWrapper() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, ctxi)
}

func CtxWithRequest(ctx context.Context, r *http.Request) *Ctx {
	ctxi := newCtx(ctx)
	ctxi.setWithReq(r)
	return ctxi
}

func CtxFromRequest(r *http.Request) *Ctx {
	ctxi := newCtx(r.Context())
	ctxi.setWithReq(r)
	return ctxi
}

func CtxFromContext(ctx context.Context) *Ctx {
	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*Ctx)
	if !ok {
		c = newCtx(ctx)
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	return c
}

func newCtx(ctx context.Context) *Ctx {
	span := trace.FromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID.String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	return &Ctx{
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

func (c *Ctx) setWithReq(r *http.Request) {
	c.Request = r
	c.Token = httpi.GetToken(r)
	c.DeviceInfo = Device(r.Header)
	c.Internal = r.Header.Get(httpi.GrpcInternal)
}

func (c *Ctx) reset(ctx context.Context) *Ctx {
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

func (c *Ctx) GetAuthInfo(auth func(*Ctx) error) (AuthInfo, error) {
	if c.Authorization == nil {
		c.Authorization = new(Authorization)
	}
	if err := auth(c); err != nil {
		return nil, err
	}
	return c.AuthInfo, nil
}

func init() {
	jwt.WithUnmarshaller(jwtUnmarshaller)(jwti.Parser)
}

func jwtUnmarshaller(ctx jwt.CodingContext, data []byte, v interface{}) error {
	if ctx.FieldDescriptor == jwt.ClaimsFieldDescriptor {
		if c, ok := (*v.(*jwt.Claims)).(*Authorization); ok {
			c.Token = stringsi.ToString(data)
			return json.Unmarshal(data, c)
		}
	}
	return json.Unmarshal(data, v)
}
