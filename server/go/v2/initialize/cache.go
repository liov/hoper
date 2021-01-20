package initialize

import (
	"github.com/liov/hoper/go/v2/utils/dao/cache"
)

func (init *Init) P3Cache() cache.Cache {
	return cache.New(1000).LRU().Build()
}
