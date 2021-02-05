package user

import (
	"context"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/net"
	"github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	jwti "github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"go.opencensus.io/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

//Cannot use 'resumes' (type []*model.Resume) as type []CmpKey
//我认为这是一个bug
//[]int可以是interface，却不可以是[]interface
//var test []array.CmpKey
//test = append(test,resumes[0]) 可行
//test = append(test,resumes...) 不可行，可笑
func (x *Resume) CmpKey() uint64 {
	return x.Id
}

var UserserviceServicedesc = &_UserService_serviceDesc

/*
func RegisterUserServiceHandlerFromModuleWithReConnect(ctx context.Context, mux *runtime.ServeMux, getEndPort func() string, opts []grpc.DialOption) (err error) {
	endPort:=getEndPort()
	conn, err := grpc.Dial(endPort, opts...)
	if err != nil {
		return err
	}
	client := NewUserServiceClient(conn)
	reconn.ReConnectMap[endPort] = reconn.ReConnect(client, getEndPort, opts)
	return RegisterUserServiceHandlerClient(ctx, mux, client)
}
*/

/*----------------------------ORM-------------------------------*/
func userTableName(id uint64) string {
	if id < 1_000_000 {
		return "user"
	}
	return "user_" + string(byte(id/1_000_000+49))
}
func (x *User) TableName() string {
	return userTableName(x.Id)
}

func (x *UserBaseInfo) TableName() string {
	return userTableName(x.Id)
}

func (x *UserAuthInfo) TableName() string {
	return userTableName(x.Id)
}

func (x *Resume) TableName() string {
	return "resume"
}

func (x *EditReq) TableName() string {
	return userTableName(x.Id)
}

func (x *AuthInfo) TableName() string {
	return userTableName(x.Id)
}

/*----------------------------CTX上下文-------------------------------*/
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

type AuthInfo struct {
	Id           uint64     `json:"id"`
	IdStr        string     `json:"-" gorm:"-"`
	Name         string     `json:"name"`
	Role         Role       `json:"role"`
	Status       UserStatus `json:"status"`
	LastActiveAt int64      `json:"lat,omitempty"`
	ExpiredAt    int64      `json:"exp,omitempty"`
	LoginAt      int64      `json:"iat,omitempty"`
}

func (x *AuthInfo) UserAuthInfo() *UserAuthInfo {
	return &UserAuthInfo{
		Id:           x.Id,
		Name:         x.Name,
		Role:         x.Role,
		Status:       x.Status,
		LastActiveAt: x.LastActiveAt,
		ExpiredAt:    x.ExpiredAt,
		LoginAt:      x.LoginAt,
	}
}

type Authorization struct {
	*AuthInfo
	Token string `json:"-"`
}

func (x *Authorization) Valid(helper *jwt.ValidationHelper) error {
	if x.ExpiredAt != 0 && x.LastActiveAt > x.ExpiredAt {
		return UserErrLoginTimeout
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
	x.IdStr = strconv.FormatUint(x.Id, 10)
	return nil
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

func (x *DeviceInfo) UserDeviceInfo() *UserDeviceInfo {
	return &UserDeviceInfo{
		Device:     x.Device,
		Os:         x.Os,
		AppCode:    x.AppCode,
		AppVersion: x.AppVersion,
		IP:         x.IP,
		Lng:        x.Lng,
		Lat:        x.Lat,
		Area:       x.Area,
		UserAgent:  x.UserAgent,
	}
}

type Ctx struct {
	context.Context `json:"-"`
	TraceID         string `json:"-"`
	Authorization
	*DeviceInfo                `json:"-"`
	RequestAt                  time.Time   `json:"-"`
	RequestUnix                int64       `json:"iat,omitempty"`
	Header                     metadata.MD `json:"-"`
	Peer                       *peer.Peer  `json:"-"`
	grpc.ServerTransportStream `json:"-"`
	Internal                   string
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

func (c *Ctx) SetHeader(md metadata.MD) error {
	for k, v := range md {
		c.Header[k] = v
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
		c.Header[k] = v
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
	c.Header[k] = []string{v}

	if c.ServerTransportStream != nil {
		err := c.ServerTransportStream.SendHeader(metadata.MD{k: []string{v}})
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Ctx) SetCookie(v string) error {
	c.Header[httpi.HeaderSetCookie] = []string{v}

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
		c.Header[k] = v
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

type ctxKey struct{}

func CtxWithRequest(ctx context.Context, r *http.Request) context.Context {
	ctxi := newCtx(ctx)
	ctxi.setWithReq(r)
	return context.WithValue(context.Background(), ctxKey{}, ctxi)
}

func ConvertContext(r *http.Request) pick.Context {
	ctxi := r.Context().Value(ctxKey{})
	c, ok := ctxi.(*Ctx)
	if !ok {
		c = newCtx(r.Context())
		c.setWithReq(r)
	}
	return c
}

func GetAuthInfo(c context.Context) string {
	return CtxFromContext(c).Token
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
		c.WithContext(grpc.NewContextWithServerTransportStream(c.Context, stream))
		c.ServerTransportStream = stream
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
		Authorization: Authorization{
			AuthInfo: new(AuthInfo),
		},
		RequestAt:   now,
		RequestUnix: now.Unix(),
		Logger:      log.Default.With(zap.String("TraceId", traceId)),
	}
}

func (c *Ctx) setWithReq(r *http.Request) {
	isGrpc := r.ProtoMajor == 2 && strings.Contains(r.Header.Get(httpi.HeaderContentType), httpi.ContentGRPCHeaderValue)
	var p *peer.Peer
	if !isGrpc {
		p = &peer.Peer{
			Addr: neti.StrAddr(r.RemoteAddr),
		}
		if r.TLS != nil {
			p.AuthInfo = credentials.TLSInfo{State: *r.TLS, CommonAuthInfo: credentials.CommonAuthInfo{SecurityLevel: credentials.PrivacyAndIntegrity}}
		}
	}
	c.Token = httpi.GetToken(r)
	c.Header = metadata.MD(r.Header)
	c.DeviceInfo = Device(r.Header)
	c.Internal = r.Header.Get(httpi.HeaderInternal)
}

func (c *Ctx) GetAuthInfo(auth func(*Ctx) error) (*AuthInfo, error) {
	if c.AuthInfo == nil {
		c.AuthInfo = new(AuthInfo)
	}
	if err := auth(c); err != nil {
		return nil, err
	}
	return c.AuthInfo, nil
}

func JWTUnmarshaller(ctx jwt.CodingContext, data []byte, v interface{}) error {
	if ctx.FieldDescriptor == jwt.ClaimsFieldDescriptor {
		if c, ok := (*v.(*jwt.Claims)).(*Authorization); ok {
			c.Token = stringsi.ToString(data)
			return json.Unmarshal(data, c.AuthInfo)
		}
	}
	return json.Unmarshal(data, v)
}
