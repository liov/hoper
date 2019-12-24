package initialize

import (
	"github.com/bluele/gcache"
)

func (init *Init) P3Cache() gcache.Cache {
	return gcache.New(20).LRU().Build()
}
