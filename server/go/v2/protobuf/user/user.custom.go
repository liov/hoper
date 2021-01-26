package user

import (
	"context"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/liov/hoper/go/v2/utils/encoding/json"
	"github.com/liov/hoper/go/v2/utils/net/http"
	"github.com/liov/hoper/go/v2/utils/net/http/pick"
	stringsi "github.com/liov/hoper/go/v2/utils/strings"
	jwti "github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"go.opencensus.io/trace"
	"google.golang.org/grpc/metadata"
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

type AuthInfo struct {
	Id           uint64     `json:"id"`
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

func (x *AuthInfo) Valid(helper *jwt.ValidationHelper) error {
	if x.ExpiredAt != 0 && x.LastActiveAt > x.ExpiredAt {
		return UserErr_LoginTimeout
	}
	return nil
}

func (x *AuthInfo) GenerateToken(secret []byte) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func (x *AuthInfo) ParseToken(req *http.Request, secret string) error {
	return jwti.ParseToken(x, httpi.GetToken(req), secret)
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
	*AuthInfo
	Authorization string `json:"-"`
	*DeviceInfo   `json:"-"`
	RequestAt     time.Time   `json:"-"`
	RequestUnix   int64       `json:"iat,omitempty"`
	MD            metadata.MD `json:"-"`
	parsed        bool
}

func (c *Ctx) StartSpan(name string, o ...trace.StartOption) (*Ctx, *trace.Span) {
	ctx, span := trace.StartSpan(c.Context, name, o...)
	c.Context = ctx
	if c.TraceID == "" {
		c.TraceID = span.SpanContext().TraceID.String()
	}
	return c, span
}

type ctxKey struct{}

func CtxWithRequest(ctx context.Context, r *http.Request) context.Context {
	span := trace.FromContext(ctx)
	now := time.Now()
	user := new(AuthInfo)
	user.LastActiveAt = now.Unix()
	return context.WithValue(context.Background(), ctxKey{},
		&Ctx{
			Context:       ctx,
			TraceID:       span.SpanContext().TraceID.String(),
			AuthInfo:      user,
			Authorization: httpi.GetToken(r),
			DeviceInfo:    Device(r.Header),
			RequestAt:     now,
			RequestUnix:   user.LastActiveAt,
		})
}

func ConvertContext(r *http.Request) pick.Claims {
	ctxi := r.Context().Value(ctxKey{})
	if c, ok := ctxi.(*Ctx); ok {
		if !c.parsed {
			c.MD = metadata.MD(r.Header)
			c.parsed = true
		}
		return c
	}
	return NewCtx(r.Context())
}

func Authorization(c context.Context) string {
	return CtxFromContext(c).Authorization
}

func CtxFromContext(ctx context.Context) *Ctx {
	ctxi := ctx.Value(ctxKey{})
	if c, ok := ctxi.(*Ctx); ok {
		if !c.parsed {
			md, _ := metadata.FromIncomingContext(ctx)
			c.MD = md
			c.parsed = true
		}
		return c
	}
	return NewCtx(ctx)
}

func NewCtx(ctx context.Context) *Ctx {
	now := time.Now()
	user := new(AuthInfo)
	user.LastActiveAt = now.Unix()
	return &Ctx{
		Context:     ctx,
		AuthInfo:    user,
		RequestAt:   now,
		RequestUnix: user.LastActiveAt,
	}
}

func (c *Ctx) GetAuthInfo(auth func(*Ctx) (*AuthInfo, error)) (*AuthInfo, error) {
	if c.AuthInfo == nil {
		c.AuthInfo = new(AuthInfo)
	}
	user, err := auth(c)
	if err != nil {
		return nil, err
	}
	c.AuthInfo = user
	return user, nil
}

func JWTUnmarshaller(ctx jwt.CodingContext, data []byte, v interface{}) error {
	if ctx.FieldDescriptor == jwt.ClaimsFieldDescriptor {
		if c, ok := (*v.(*jwt.Claims)).(*Ctx); ok {
			c.Authorization = stringsi.ToString(data)
			return json.Unmarshal(data, c.AuthInfo)
		}
	}
	return json.Unmarshal(data, v)
}
