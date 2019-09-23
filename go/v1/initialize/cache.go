package initialize

import (
	"github.com/bluele/gcache"
	"github.com/liov/hoper/go/v1/initialize/dao"
)

func (i *Init) P3Cache() {
	dao.SetCache(gcache.New(20).LRU().Build())
}
