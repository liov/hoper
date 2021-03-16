package initialize

import (
	"github.com/dgraph-io/ristretto"
)
// 考虑换cache，ristretto存一个值，循环取居然还会miss,某个issue提要内存占用过大，直接初始化1.5MB
// freecache不能存对象，可能要为每个对象写UnmarshalBinary 和 MarshalBinary
// go-cache
func (init *Init) P3Cache() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 27, // maximum cost of cache (125MB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return cache
}
