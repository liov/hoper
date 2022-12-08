package download

import "github.com/liov/hoper/server/go/lib/utils/generics/conctrl"

const (
	KindGetAllPhoto conctrl.Kind = iota
	KindGetPhoto
	KindDownloadPhoto
	KindGetFollow
	KindGetAllFollow
)

const (
	Referer = "https://m.weibo.cn/"
)
