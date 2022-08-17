package context

import (
	"time"
)

// 不现实，context里可能装任何数据，泛型会限定只能装一种，除非限定就是any
type Context[K any, V any] interface {
	Deadline() (deadline time.Time, ok bool)
	Done() <-chan struct{}
	Err() error
	Value(key K) V
}
