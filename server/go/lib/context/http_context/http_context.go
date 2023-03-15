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
