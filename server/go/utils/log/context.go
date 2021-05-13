package log

import (
	"context"

	"github.com/google/uuid"
)

type IdKey struct{}

// NewContext returns a new Context that carries value u.
func NewContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, IdKey{}, uuid.New().String())
}

// FromContext returns the User value stored in ctx, if any.
func FromContext(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(IdKey{}).(string)
	return u, ok
}
