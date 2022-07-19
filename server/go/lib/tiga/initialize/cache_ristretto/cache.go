package cache_ristretto

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/dgraph-io/ristretto"
)

type CacheConfig struct {
	NumCounters        int64
	MaxCost            int64
	BufferItems        int64
	Metrics            bool
	IgnoreInternalCost bool
}

func (conf *CacheConfig) generate() *ristretto.Cache {
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

func (conf *CacheConfig) Generate() interface{} {
	return conf.generate()
}

// 考虑换cache，ristretto存一个值，循环取居然还会miss(没开IgnoreInternalCost的原因),某个issue提要内存占用过大，直接初始化1.5MB
// freecache不能存对象，可能要为每个对象写UnmarshalBinary 和 MarshalBinary
// go-cache

type Cache struct {
	*ristretto.Cache
	Conf CacheConfig
}

func (c *Cache) Config() initialize.Generate {
	return &c.Conf
}

func (c *Cache) SetEntity(entity interface{}) {
	if cache, ok := entity.(*ristretto.Cache); ok {
		c.Cache = cache
	}
}

func (e *Cache) Close() error {
	e.Cache.Close()
	return nil
}
