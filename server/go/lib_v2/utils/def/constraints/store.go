package constraints

import "time"

type StoreWithExpire[K, V any] interface {
	Set(k K, v V, expire time.Duration)
	Get(k K) V
}

type Store[K, V any] interface {
	Set(k K, v V)
	Get(k K) V
}
