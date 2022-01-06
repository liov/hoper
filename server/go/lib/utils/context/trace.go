package contexti

import (
	"context"
	logi "github.com/actliboy/hoper/server/go/lib/utils/log"
	"go.uber.org/zap"
)

func TraceId(ctx context.Context) zap.Field {
	if traceId, ok := ctx.Value(traceIdKey{}).(string); ok {
		zap.String(logi.TraceId, traceId)
	}
	return zap.String(logi.TraceId, "unknown")
}

type traceIdKey struct{}

func SetTranceId(traceId string) context.Context {
	return context.WithValue(context.Background(), traceIdKey{}, traceId)
}
