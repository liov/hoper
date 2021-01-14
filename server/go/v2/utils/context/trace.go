package contexti

import (
	"context"

	"github.com/google/uuid"
)

type traceKey struct{}

func TraceCtx() context.Context {
	return context.WithValue(context.Background(), traceKey{}, uuid.New())
}

func GetTraceId(ctx context.Context) uuid.UUID {
	traceId, ok := ctx.Value(traceKey{}).(uuid.UUID)
	if ok{
		return traceId
	}
	return uuid.New()
}