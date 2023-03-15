package ristretto

import (
	"github.com/dgraph-io/ristretto"
)

type CacheConfig struct {
	NumCounters        int64
	MaxCost            int64
	BufferItems        int64
	Metrics            bool
	IgnoreInternalCost bool
}

func (conf *CacheConfig) Build() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters:        conf.NumCounters, // number of keys to track frequency of (10M).
		MaxCost:            conf.MaxCost,     // maximum cost of cache (MaxCost * 1MB).
		BufferItems:        64,               // number of keys per Get buffer.
		Metrics:            conf.Metrics,     // number of keys per Get buffer.
		IgnoreInternalCost: conf.IgnoreInternalCost,
	})
	if err != nil {
		panic(err)
	}
	return cache
}

// 考虑换cache，ristretto存一个值，循环取居然还会miss(没开IgnoreInternalCost的原因),某个issue提要内存占用过大，直接初始化1.5MB
// freecache不能存对象，可能要为每个对象写UnmarshalBinary 和 MarshalBinary
// go-cache

type Cache struct {
	*ristretto.Cache
	Conf CacheConfig
}

func (c *Cache) Config() any {
	return &c.Conf
}

func (c *Cache) SetEntity() {
	c.Cache = c.Conf.Build()
}

func (e *Cache) Close() error {
	e.Cache.Close()
	return nil
}
