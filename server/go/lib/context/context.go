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

type Ctx struct {
	LastActiveAt int64
	*Authorization
	*contexti.RequestContext
}

func (c *Ctx) StartSpan(name string, o ...trace.StartOption) (*Ctx, *trace.Span) {
	_, span := c.RequestContext.StartSpan(name, o...)
	return c, span
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
		return &Ctx{Authorization: &Authorization{}, RequestContext: contexti.NewCtx(ctx)}
	}
	if c.ServerTransportStream == nil {
		c.ServerTransportStream = grpc.ServerTransportStreamFromContext(ctx)
	}
	return c
}
