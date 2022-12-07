package conf

import "github.com/liov/hoper/server/go/lib/utils/fs"

type serverConfig struct {
	PassSalt    string
	TokenMaxAge int64
	TokenSecret string
	PageSize    int8

	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSaveDir fs.Dir //二维码保存路径
	PrefixUrl     string
	FontSaveDir   fs.Dir //字体保存路径

	CrawlerName string //爬虫

	Moment Moment
}

type Limit struct {
	SecondLimitKey, MinuteLimitKey, DayLimitKey       string
	SecondLimitCount, MinuteLimitCount, DayLimitCount int64
}

type Moment struct {
	MaxContentLen int
	Limit         Limit
}
