package contexti

import (
	"context"
	"errors"
	"github.com/liov/hoper/go/v2/protobuf/utils/errorcode"
	"github.com/liov/hoper/go/v2/utils/net/http/request"
	timei "github.com/liov/hoper/go/v2/utils/time"
	"google.golang.org/grpc/metadata"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/log"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	jwti "github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
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
	unknow := true
	var info DeviceInfo
	//Device-Info:device,osInfo,appCode,appVersion
	if infos := r.Values(httpi.HeaderDeviceInfo); len(infos) == 4 {
		unknow = false
		info.Device = infos[0]
		info.Os = infos[1]
		info.AppCode = infos[2]
		info.AppVersion = infos[3]

	}
	// area:xxx
	// location:1.23456,2.123456
	if area := r.Get(httpi.HeaderArea); area != "" {
		unknow = false
		info.Area, _ = url.PathUnescape(area)
	}
	if infos := r.Values(httpi.HeaderLocation); len(infos) == 2 {
		unknow = false
		info.Lng = infos[0]
		info.Lat = infos[1]

	}

	if userAgent := r.Get(httpi.HeaderUserAgent); userAgent != "" {
		unknow = false
		info.UserAgent = userAgent
	}
	if ip := r.Get(httpi.HeaderXForwardedFor); ip != "" {
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
	Request *http.Request
	grpc.ServerTransportStream
	Internal string
	*log.Logger
}

var _ = pick.Context(new(Ctx))

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

func CtxWithRequest(ctx context.Context, r *http.Request) context.Context {
	ctxi := newCtx(ctx)
	ctxi.setWithReq(r)
	return context.WithValue(context.Background(), ctxKey{}, ctxi)
}

func ConvertContext(r *http.Request) *Ctx {
	ctxi := r.Context().Value(ctxKey{})
	c, ok := ctxi.(*Ctx)
	if !ok {
		c = newCtx(r.Context())
		c.setWithReq(r)
	}
	return c
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
		// 每个请求对应一个实例，后续并发量大考虑移除直接使用log库实例
		Logger: log.Default.With(zap.String("traceId", traceId)),
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
	c.Logger = log.Default.With(zap.String("traceId", traceId))
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

func (c *Ctx) SetHeader(md metadata.MD) error {
	for k, v := range md {
		c.Request.Header[k] = v
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Ctx) SendHeader(md metadata.MD) error {
	for k, v := range md {
		c.Request.Header[k] = v
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Ctx) WriteHeader(k, v string) error {
	c.Request.Header[k] = []string{v}

	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Ctx) SetCookie(v string) error {
	c.Request.Header[httpi.HeaderSetCookie] = []string{v}

	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{httpi.HeaderSetCookie: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Ctx) SetTrailer(md metadata.MD) error {
	for k, v := range md {
		c.Request.Header[k] = v
	}
	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SetTrailer(md)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Ctx) Method() string {
	if c.ServerTransportStream != nil {
		return c.ServerTransportStream.Method()
	}
	return ""
}

func (c *Ctx) HandleError(err error) {
	if err != nil {
		c.Error(err.Error())
	}
}

func (c *Ctx) ErrorLog(err, originErr error, funcName string) error {
	// caller 用原始logger skip刚好
	c.Logger.Logger.Error(originErr.Error(), zap.Int("type", errorcode.Code(err)), zap.String(log.Position, funcName))
	return err
}

type AuthInfoDao struct {
	Secret    string
	AuthCache *ristretto.Cache
	AuthPool  *sync.Pool
	Update    func(*Ctx) error
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

func (c *Ctx) GeToken() string {
	return c.Token
}

func (c *Ctx) GetReqAt() *request.RequestAt {
	return &c.RequestAt
}

func (c *Ctx) GetLogger() *log.Logger {
	return c.Logger
}
