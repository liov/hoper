package initialize

import (
	"github.com/dgraph-io/ristretto"
	reflecti "github.com/liov/hoper/go/v2/utils/reflect"
)

type CacheConfig struct {
	NumCounters        int64
	MaxCost            int64
	BufferItems        int64
	Metrics            bool
	IgnoreInternalCost bool
}

func (conf *CacheConfig) Generate() *ristretto.Cache {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: conf.NumCounters, // number of keys to track frequency of (10M).
		MaxCost:     conf.MaxCost,     // maximum cost of cache (125MB).
		BufferItems: 64, // number of keys per Get buffer.
		Metrics:     conf.Metrics,     // number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
	return cache
}

// 考虑换cache，ristretto存一个值，循环取居然还会miss,某个issue提要内存占用过大，直接初始化1.5MB
// freecache不能存对象，可能要为每个对象写UnmarshalBinary 和 MarshalBinary
// go-cache
func (init *Init) P3Cache() *ristretto.Cache {
	conf := &CacheConfig{}
	if exist := reflecti.GetFieldValue(init.conf, conf); !exist {
		return nil
	}

	return conf.Generate()
}

/*func (init *Init) P3Cache() *cache1.Cache {
	return cache1.New(5*time.Minute, 10*time.Minute)
}*/
