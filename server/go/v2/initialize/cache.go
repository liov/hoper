package initialize

import (
	"github.com/liov/hoper/go/v2/utils/cache"
)

func (init *Init) P3Cache() cache.Cache {
	return cache.New(20).LRU().Build()
}
