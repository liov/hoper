package conf

import "github.com/liov/hoper/go/v2/utils/fs"

type serverConfig struct {
	PassSalt    string
	TokenMaxAge int64
	TokenSecret string
	PageSize    int8

	UploadDir      fs.Dir
	UploadMaxSize  int64
	UploadAllowExt []string

	LogSaveDir  fs.Dir
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSaveDir fs.Dir //二维码保存路径
	PrefixUrl     string
	FontSaveDir   fs.Dir //字体保存路径

	CrawlerName string //爬虫
	Limit Limit
}

type Limit struct {
	SecondLimit, MinuteLimit, DayLimit                string
	SecondLimitCount, MinuteLimitCount, DayLimitCount int64
}

