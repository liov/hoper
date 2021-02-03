package auth

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go/v4"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	neti "github.com/liov/hoper/go/v2/utils/net"
	httpi "github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	jwti "github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Info interface {
	jwt.Claims
	GenerateToken(secret []byte) (string, error)
	ParseToken(ctx *Ctx, secret string) error
}

type UserInfo struct {
	Id           uint64     `json:"id"`
	IdStr        string     `json:"-" gorm:"-"`
	Name         string     `json:"name"`
	Role         int8       `json:"role"`
	Status       int8 `json:"status"`
	LastActiveAt int64      `json:"lat,omitempty"`
	ExpiredAt    int64      `json:"exp,omitempty"`
	LoginAt      int64      `json:"iat,omitempty"`
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

func (x *UserInfo) ParseToken(token, secret string) error {
	if err := jwti.ParseToken(x, token, secret); err != nil {
		return err
	}
	x.IdStr = strconv.FormatUint(x.Id, 10)
	return nil
}

type Cache struct {
	Info
	Authorization string
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
	//Device-Info:device-osInfo-appCode-appVersion
	if deviceInfo := r.Get(httpi.HeaderDeviceInfo); deviceInfo != "" {
		unknow = false
		infos := strings.Split(deviceInfo, "-")
		if len(infos) == 4 {
			info.Device = infos[0]
			info.Os = infos[1]
			info.AppCode = infos[2]
			info.AppVersion = infos[3]
		}
	}
	// area:xxx
	// location:1.23456,2.123456
	if area := r.Get(httpi.HeaderArea); area != "" {
		unknow = false
		info.Area, _ = url.PathUnescape(area)
	}
	if location := r.Get(httpi.HeaderLocation); location != "" {
		unknow = false
		infos := strings.Split(location, ",")
		if len(infos) == 2 {
			info.Lng = infos[0]
			info.Lat = infos[1]
		}
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
	context.Context `json:"-"`
	TraceID         string `json:"-"`
	Info
	Authorization              string `json:"-"`
	*DeviceInfo                `json:"-"`
	RequestAt                  time.Time   `json:"-"`
	RequestUnix                int64       `json:"iat,omitempty"`
	Header                     metadata.MD `json:"-"`
	Peer                       *peer.Peer  `json:"-"`
	grpc.ServerTransportStream `json:"-"`
	grpc                       bool
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
	span := trace.FromContext(ctx)
	now := time.Now()
	return context.WithValue(context.Background(), ctxKey{},
		&Ctx{
			Context:       ctx,
			TraceID:       span.SpanContext().TraceID.String(),
			Authorization: httpi.GetToken(r),
			DeviceInfo:    Device(r.Header),
			RequestAt:     now,
			RequestUnix:   now.Unix(),
			Header:        metadata.MD(r.Header),
		})
}

func ConvertContext(r *http.Request) pick.Context {
	ctxi := r.Context().Value(ctxKey{})
	c, ok := ctxi.(*Ctx)
	if !ok {
		c = newCtx(r.Context())
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get(httpi.HeaderContentType), httpi.ContentGRPCHeaderValue) {
			c.grpc = true
		}
	}
	if c.Header == nil {
		c.Header = metadata.MD(r.Header)
	}
	if c.Peer == nil && !c.grpc {
		p := &peer.Peer{
			Addr: neti.StrAddr(r.RemoteAddr),
		}
		if r.TLS != nil {
			p.AuthInfo = credentials.TLSInfo{State: *r.TLS, CommonAuthInfo: credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity}}
		}
		c.Peer = p
	}
	if c.ServerTransportStream == nil && !c.grpc {
		c.ServerTransportStream = new(runtime.ServerTransportStream)
	}
	return c
}

func Authorization(c context.Context) string {
	return CtxFromContext(c).Authorization
}

func CtxFromContext(ctx context.Context) *Ctx {
	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*Ctx)
	if !ok {
		c = newCtx(ctx)
	}
	if c.Header == nil {
		md, _ := metadata.FromIncomingContext(ctx)
		c.Header = md
	}
	if c.Peer == nil {
		p, _ := peer.FromContext(ctx)
		c.Peer = p
	}
	if c.ServerTransportStream == nil {
		var stream grpc.ServerTransportStream
		if stream = grpc.ServerTransportStreamFromContext(ctx); stream == nil {
			stream = new(runtime.ServerTransportStream)
		}
	}
	return c
}

func newCtx(ctx context.Context) *Ctx {
	now := time.Now()

	return &Ctx{
		Context:     ctx,
		RequestAt:   now,
		RequestUnix: now.Unix(),
	}
}

func (c *Ctx) GetAuthInfo(auth func(*Ctx) error) (Info, error) {
	if err := auth(c); err != nil {
		return nil, err
	}
	return c.Info, nil
}

func init() {
	jwt.WithUnmarshaller(JWTUnmarshaller)(jwti.Parser)
}


func JWTUnmarshaller(ctx jwt.CodingContext, data []byte, v interface{}) error {
	if ctx.FieldDescriptor == jwt.ClaimsFieldDescriptor {
		if c, ok := (*v.(*jwt.Claims)).(*Ctx); ok {
			c.Authorization = stringsi.ToString(data)
			return json.Unmarshal(data, c.Info)
		}
	}
	return json.Unmarshal(data, v)
}
