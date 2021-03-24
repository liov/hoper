package koko

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strings"
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
	jwt.Claims
	IdStr() string
}

type UserInfo struct {
	Id           uint64 `json:"id"`
	IdStr        string `json:"-" gorm:"-"`
	Name         string `json:"name"`
	Role         int8   `json:"role"`
	Status       int8   `json:"status"`
	LastActiveAt int64  `json:"lat,omitempty"`
	ExpiredAt    int64  `json:"exp,omitempty"`
	LoginAt      int64  `json:"iat,omitempty"`
}

func (x *UserInfo) Valid(helper *jwt.ValidationHelper) error {
	if x.ExpiredAt != 0 && x.LastActiveAt > x.ExpiredAt {
		return errors.New("登录超时")
	}
	return nil
}

func (x *UserInfo) GenerateToken(secret []byte) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

type Authorization struct {
	AuthInfo
	IdStr string
	Token string
}

type DeviceInfo struct {
	//设备
	Device     string `json:"device" gorm:"size:255"`
	Os         string `json:"os" gorm:"size:255"`
	AppCode    string `json:"appCode" gorm:"size:255"`
	AppVersion string `json:"appVersion" gorm:"size:255"`
	IP         string `json:"IP" gorm:"size:255"`
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
	RequestAt   time.Time
	RequestUnix int64
	Request     *http.Request
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
		Context: ctx,
		TraceID: traceId,
		Authorization: &Authorization{},
		RequestAt:   now,
		RequestUnix: now.Unix(),
		Logger:      log.Default.With(zap.String("traceId", traceId)),
	}
}

func (c *Ctx) setWithReq(r *http.Request) {
	c.Request = r
	c.Token = httpi.GetToken(r)
	c.DeviceInfo = Device(r.Header)
	c.Internal = r.Header.Get(httpi.GrpcInternal)
}


func (c *Ctx) GeToken() string {
	return c.Token
}

func (c *Ctx) GetReqTime() time.Time {
	return c.RequestAt
}

func (c *Ctx) GetLogger() *log.Logger {
	return c.Logger
}

func (c *Ctx) reset(ctx context.Context) *Ctx {
	span := trace.FromContext(ctx)
	now := time.Now()
	traceId := span.SpanContext().TraceID.String()
	if traceId == "" {
		traceId = uuid.New().String()
	}
	c.Context = ctx
	c.RequestAt = now
	c.RequestUnix = now.Unix()
	c.Logger = log.Default.With(zap.String("traceId", traceId))
	return c
}

type AuthInfoDao struct {
	Secret    string
	AuthCache *ristretto.Cache
	AuthPool  *sync.Pool
	Update    func(*Ctx) error
}

func (c *Ctx) GetAuthInfo(authDao AuthInfoDao,update bool) error {
	var signature string
	if authDao.AuthCache != nil {
		signature = c.Token[strings.LastIndexByte(c.Token, '.')+1:]
		cacheTmp, ok := authDao.AuthCache.Get(signature)
		if ok {
			if authorization, ok := cacheTmp.(*Authorization); ok {
				c.Authorization = authorization
				return nil
			}
		}
	}
	if authDao.Secret != "" {
		c.AuthInfo = authDao.AuthPool.Get().(AuthInfo)
		if err := jwti.ParseToken(c.Authorization, c.Token, authDao.Secret); err != nil {
			return err
		}
	}
	if update && authDao.Update != nil {
		if err := authDao.Update(c); err != nil {
			return err
		}
	}
	c.IdStr = c.AuthInfo.IdStr()
	if authDao.AuthCache != nil {
		authDao.AuthCache.SetWithTTL(signature, c.Authorization,0, 5*time.Second)
	}
	return nil
}

func init() {
	jwt.WithUnmarshaller(jwtUnmarshaller)(jwti.Parser)
}

func jwtUnmarshaller(ctx jwt.CodingContext, data []byte, v interface{}) error {
	if ctx.FieldDescriptor == jwt.ClaimsFieldDescriptor {
		if c, ok := (*v.(*jwt.Claims)).(*Authorization); ok {
			c.Token = stringsi.ToString(data)
			return json.Unmarshal(data, c.AuthInfo)
		}
	}
	return json.Unmarshal(data, v)
}
