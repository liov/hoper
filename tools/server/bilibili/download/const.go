package download

import "github.com/actliboy/hoper/server/go/lib/utils/conctrl"

const (
	KindGetFavListUrl conctrl.Kind = 0
	KindViewInfo      conctrl.Kind = 1
	KindDownloadCover conctrl.Kind = 2
	KindGetPlayerUrl  conctrl.Kind = 3
	KindDownloadVideo conctrl.Kind = 4
)

const (
	VideoTypeFlv        = 0
	VideoTypeM4sCodec12 = 1
	VideoTypeM4sCodec7  = 2
)
