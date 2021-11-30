package contexti

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	contexti "github.com/liov/hoper/server/go/lib/utils/context"
	"github.com/liov/hoper/server/go/lib/utils/encoding/json"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	jwti "github.com/liov/hoper/server/go/lib/utils/verification/auth/jwt"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"net/http"
	"sync"
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
	AuthInfoRaw  string `json:"-"`
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

type Ctx struct {
	*Authorization
	*contexti.RequestContext
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

func CtxFromRequest(r *http.Request, tracing bool) (*Ctx, *trace.Span) {
	ctxi, span := contexti.CtxWithRequest(r, tracing)
	return &Ctx{Authorization: &Authorization{}, RequestContext: ctxi}, span
}

type ctxKey struct{}

func (ctxi *Ctx) ContextWrapper() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, ctxi)
}

func CtxFromContext(ctx context.Context) *Ctx {
	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*Ctx)
	if !ok {
		ctxi := contexti.NewCtx(ctx)
		c = &Ctx{Authorization: &Authorization{}, RequestContext: ctxi}
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	return c
}

func init() {
	jwt.WithUnmarshaller(jwtUnmarshaller)(jwti.Parser)
}

func jwtUnmarshaller(ctx jwt.CodingContext, data []byte, v interface{}) error {
	if ctx.FieldDescriptor == jwt.ClaimsFieldDescriptor {
		if c, ok := (*v.(*jwt.Claims)).(*Authorization); ok {
			c.AuthInfoRaw = stringsi.ToString(data)
			return json.Unmarshal(data, c)
		}
	}
	return json.Unmarshal(data, v)
}
