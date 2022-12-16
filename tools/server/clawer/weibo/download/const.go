package download

import (
	"github.com/liov/hoper/server/go/lib_v2/utils/conctrl"
	"time"
)

const (
	KindNormal conctrl.Kind = iota
	KindGet
	KindDownload
)

const (
	Referer = "https://m.weibo.cn/"
)

const (
	TimeFormat = time.RubyDate
)
