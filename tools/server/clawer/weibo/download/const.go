package download

import (
	"github.com/liov/hoper/server/go/lib/v2/utils/conctrl"
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
