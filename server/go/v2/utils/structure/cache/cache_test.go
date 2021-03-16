package cache

import (
	"fmt"
	"runtime/debug"
	"sync"
	"testing"
	"time"

	"github.com/dgraph-io/ristretto"
	"github.com/liov/hoper/go/v2/utils/structure/cache/freecache"
	"github.com/liov/hoper/go/v2/utils/structure/cache/gcache"
	cache1 "github.com/liov/hoper/go/v2/utils/structure/cache/go-cache"
)

func BenchmarkFree(b *testing.B){
	cacheSize := 100 * 1024 * 1024
	cache := freecache.NewCache(cacheSize)
	debug.SetGCPercent(20)
	key := []byte("abc")
	val := []byte("def")
	expire := 60 // expire in 60 seconds
	cache.Set(key, val, expire)
	for i := 0; i < b.N; i++ {

		_, err := cache.Get(key)
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGCache(b *testing.B){
	gc := gcache.New(20).
		LRU().
		Build()
	gc.Set("key", "ok")
	for i := 0; i < b.N; i++ {

		_, err := gc.Get("key")
		if err != nil {
			panic(err)
		}
	}
}

func BenchmarkGoCache(b *testing.B){
	c := cache1.New(5*time.Minute, 10*time.Minute)
	c.Set("foo", "bar", cache1.DefaultExpiration)
	for i := 0; i < b.N; i++ {

		_, found := c.Get("foo")
		if !found {
			fmt.Println("go-cache")
		}
	}
}

func BenchmarkRistretto(b *testing.B){
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
		Metrics:true,
	})
	if err != nil {
		panic(err)
	}
	// set a value with a cost of 1
	cache.Set("key", "value", 1)
	for i := 0; i < b.N; i++ {
		value, found := cache.Get("key")
		if !found {
			fmt.Println(value)
		}
	}
	fmt.Println(cache.Metrics.Ratio())
}

func BenchmarkSyncMap(b *testing.B){
	m :=sync.Map{}
	m.Store("key", "value")
	for i := 0; i < b.N; i++ {
		value, found := m.Load("key")
		if !found {
			fmt.Println(value)
		}
	}
}