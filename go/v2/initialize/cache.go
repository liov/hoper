package initialize

import (
	"github.com/bluele/gcache"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

func (i *Init) P3Cache() {
	reflect3.SetFieldValue(i.dao,gcache.New(20).LRU().Build())
}
