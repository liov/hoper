package _interface

import "golang.org/x/exp/constraints"

type Less func(a, b interface{}) bool
type Equal func(a, b interface{}) bool

type CmpKey[T constraints.Ordered] interface {
	CmpKey() T
}

type CompareField[T constraints.Ordered] interface {
	CompareField() T
}
