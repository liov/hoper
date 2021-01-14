package contexti

import (
	"context"
	"time"
)

type timeKey struct{}

// AuthContext returns a new Context that carries value u.
func AuthContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, timeKey{}, time.Now())
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (time.Time, bool) {
	now, ok := ctx.Value(timeKey{}).(time.Time)
	return now, ok
}

