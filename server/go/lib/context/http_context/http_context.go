package http_context

import (
	"context"
	contexti "github.com/liov/hoper/server/go/lib/context"
	contexti2 "github.com/liov/hoper/server/go/lib/utils/context"
	"go.opencensus.io/trace"
	"net/http"
)

type Ctx = contexti.Ctx[http.Request, uint8]

func CtxFromRequest(r *http.Request, tracing bool) (*Ctx, *trace.Span) {
	ctxi, span := contexti2.CtxWithRequest[uint8](r, tracing)
	return &Ctx{Authorization: &contexti.Authorization{}, RequestContext: ctxi}, span
}

func CtxFromContext(ctx context.Context) *Ctx {
	return contexti.CtxFromContext[http.Request, uint8](ctx)
}

func CtxContextFromRequest(r *http.Request, tracing bool) (*Context, *trace.Span) {
	ctx, span := contexti2.CtxWithRequest[*ExtProp](r, tracing)
	if ctx.Props == nil {
		ctx.Props = new(ExtProp)
	}
	ctx.Props.Init()
	return ctx, span
}

func ContextFromContext(ctx context.Context) *Context {
	ctxi := contexti2.CtxFromContext[http.Request, *ExtProp](ctx)
	if ctxi.Props == nil {
		ctxi.Props = new(ExtProp)
	}
	ctxi.Props.Init()
	return ctxi
}

type ExtProp struct {
	*contexti.Authorization
	LastActiveAt int64
}

func (e *ExtProp) Init() {
	if e.Authorization == nil {
		e.Authorization = new(contexti.Authorization)
	}
}

type Context = contexti2.RequestContext[http.Request, *ExtProp]
