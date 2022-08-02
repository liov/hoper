package constraints

import (
	"golang.org/x/exp/constraints"
	"time"
)

type Integer interface {
}

type Number interface {
}

type Ordered interface {
	constraints.Ordered | time.Time
}

type Callback[T any] interface {
	~func() | ~func() error | ~func(T) | ~func(T) error
}

type ID interface {
	constraints.Integer | ~string | ~[]byte | ~[8]byte | ~[16]byte
}
