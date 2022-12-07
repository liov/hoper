package contexti

import (
	"context"
	contexti "github.com/liov/hoper/server/go/lib/utils/context"
	httpi "github.com/liov/hoper/server/go/lib/utils/net/http"
	stringsi "github.com/liov/hoper/server/go/lib/utils/strings"
	"github.com/valyala/fasthttp"
	"net/http"
)

type ContextSource interface {
	GetContext() context.Context
	GetToken() string
}

type httpRequest http.Request

func (r *httpRequest) GetContext() context.Context {
	return (*http.Request)(r).Context()
}

func (r *httpRequest) GetToken() string {
	return r.Header.Get(httpi.HeaderAuthorization)
}

type fasthttpRequest fasthttp.Request

func (r *fasthttpRequest) GetContext() context.Context {
	return context.Background()
}

func (r *fasthttpRequest) GetToken() string {
	return stringsi.ToString(r.Header.Peek(httpi.HeaderAuthorization))
}

// TODO
func NewCtx(source ContextSource) *Ctx {
	return &Ctx{
		RequestContext: contexti.NewCtx(source.GetContext()),
	}
}
