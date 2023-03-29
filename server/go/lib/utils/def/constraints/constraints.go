package constraints

import (
	"golang.org/x/exp/constraints"
)

type Integer interface {
}

type Number interface {
	constraints.Integer | constraints.Float
}

type Callback[T any] interface {
	~func() | ~func() error | ~func(T) | ~func(T) error
}

type ID interface {
	constraints.Integer | ~string | ~[]byte | ~[8]byte | ~[16]byte
}

type Basic struct {
}
