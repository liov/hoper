package service

import (
	"context"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	model "github.com/liov/hoper/go/v2/protobuf/user"
	"github.com/liov/hoper/go/v2/user/conf"
	"github.com/liov/hoper/go/v2/user/dao"
	fasthttpi "github.com/liov/hoper/go/v2/utils/net/fasthttp"
	"github.com/liov/hoper/go/v2/utils/net/http"
	jwti "github.com/liov/hoper/go/v2/utils/verification/auth/jwt"
	"github.com/valyala/fasthttp"
	"google.golang.org/grpc/metadata"

	"github.com/liov/hoper/go/v2/utils/log"
)

func Auth(token string) (*model.UserAuthInfo, error) {
	claims := new(jwti.Claims)
	if err := jwti.ParseToken(claims, token, conf.Conf.Customize.TokenSecret); err != nil {
		return nil, err
	}
	conn := dao.NewUserRedis()
	defer conn.Close()
	user, err := conn.EfficientUserHashFromRedis(claims.UserId)
	if err != nil {
		log.Error(err)
		return nil, model.UserErr_InvalidToken
	}
	return user, nil
}

type authKey struct{}

// AuthContext returns a new Context that carries value u.
func AuthContextF(r *fasthttp.Request) context.Context {
	user, _ := Auth(fasthttpi.GetToken(r))
	return context.WithValue(context.Background(), authKey{}, user)
}

// FromContext returns the User value stored in ctx, if any.
func FromContextF(ctx context.Context) (*model.UserAuthInfo, bool) {
	user, ok := ctx.Value(authKey{}).(*model.UserAuthInfo)
	return user, ok
}

type Ctx struct {
	context.Context `json:"-"`
	TraceID         string `json:"-"`
	*model.UserAuthInfo
	authorization         string
	*model.UserDeviceInfo `json:"-"`
	RequestAt             time.Time   `json:"-"`
	RequestUnix           int64       `json:"iat,omitempty"`
	MD                    metadata.MD `json:"-"`
	parse                 bool
}

func (c *Ctx) Valid() error {
	if c.ExpiredAt != 0 && c.RequestAt.Unix() <= c.ExpiredAt {
		return model.UserErr_LoginTimeout
	}
	return nil
}

func (c *Ctx) GenerateToken() (string, error) {
	//claims.StandardClaims = jwti.NewStandardClaims(conf.Conf.Customize.TokenMaxAge, conf.Conf.Server.Domain)

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	token, err := tokenClaims.SignedString(conf.Conf.Customize.TokenSecret)
	return token, err
}

func (c *Ctx) ParseToken(req *http.Request) error {
	return jwti.ParseToken(c, httpi.GetToken(req), conf.Conf.Customize.TokenSecret)
}

func CtxFromRequest(r *http.Request) context.Context {
	return &Ctx{
		Context:        r.Context(),
		TraceID:        "",
		UserAuthInfo:   nil,
		UserDeviceInfo: model.Device(r.Header),
	}
}

func CtxWithRequest(ctx context.Context, r *http.Request) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	user, _ := Auth(httpi.GetToken(r))
	now := time.Now()
	return &Ctx{
		Context:        ctx,
		TraceID:        "",
		UserAuthInfo:   user,
		UserDeviceInfo: model.Device(r.Header),
		MD:             md,
		RequestAt:      now,
		RequestUnix:    now.Unix(),
	}
}

func CtxFromContext(ctx context.Context) *Ctx {
	if ctx, ok := ctx.(*Ctx); ok {
		return ctx
	}
	return &Ctx{
		Context:      ctx,
		TraceID:      "",
		UserAuthInfo: nil,
	}
}

func (c *Ctx) GetAuthInfo() (*model.UserAuthInfo, error) {
	user, err := Auth(c.authorization)
	if err != nil {
		return nil, err
	}
	c.UserAuthInfo = user
	return user, nil
}
