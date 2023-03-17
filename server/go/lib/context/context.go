package contexti

import (
	"context"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	contexti "github.com/liov/hoper/server/go/lib/utils/context"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	jwti "github.com/liov/hoper/server/go/lib/utils/verification/auth/jwt"
	"go.opencensus.io/trace"
	"google.golang.org/grpc"
	"sync"
)

type AuthInterface interface {
	Validate() error
	GenerateToken(secret []byte) (string, error)
	ParseToken(token string, secret []byte) error
}

func GetPool[REQ any]() sync.Pool {
	return sync.Pool{New: func() any {
		return new(Context[REQ])
	}}
}

type AuthInfo interface {
	IdStr() string
}

type Authorization struct {
	AuthInfo `json:"auth"`
	jwt.RegisteredClaims
	AuthInfoRaw string `json:"-"`
}

func (x *Authorization) Validate() error {
	return nil
}

func (x *Authorization) GenerateToken(secret []byte) (string, error) {
	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, x)
	token, err := tokenClaims.SignedString(secret)
	return token, err
}

func (x *Authorization) ParseToken(token string, secret []byte) error {
	_, err := jwti.ParseToken(x, token, secret)
	if err != nil {
		return err
	}
	x.ID = x.AuthInfo.IdStr()
	authBytes, _ := json.Marshal(x.AuthInfo)
	x.AuthInfoRaw = stringsi.ToString(authBytes)
	return nil
}

type Context[REQ any] struct {
	LastActiveAt int64
	*Authorization
	*contexti.RequestContext[REQ]
}

func (c *Context[REQ]) StartSpan(name string, o ...trace.StartOption) (*Context[REQ], *trace.Span) {
	_, span := c.RequestContext.StartSpan(name, o...)
	return c, span
}

type ctxKey struct{}

func (ctxi *Context[REQ]) ContextWrapper() context.Context {
	return context.WithValue(context.Background(), ctxKey{}, ctxi)
}

func CtxFromContext[REQ any](ctx context.Context) *Context[REQ] {
	ctxi := ctx.Value(ctxKey{})
	c, ok := ctxi.(*Context[REQ])
	if !ok {
		return &Context[REQ]{Authorization: &Authorization{}, RequestContext: contexti.NewCtx[REQ](ctx)}
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	return c
}
