package download

import "github.com/liov/hoper/server/go/lib/utils/generics/conctrl"

const (
	KindGetFavListUrl conctrl.Kind = iota
	KindViewInfo
	KindDownloadCover
	KindGetPlayerUrl
	KindDownloadVideo
	KindRecordFavList
)

const (
	VideoTypeFlv        = 0
	VideoTypeM4sCodec12 = 1
	VideoTypeM4sCodec7  = 2
)

const (
	DownloadingExt = ".downloading"
	PartEqTitle    = "!part=title!"
)