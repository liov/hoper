package conf

import (
	"runtime"
	"time"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/fs"
)

type serverConfig struct {
	Volume fs.Dir

	PassSalt string
	// 天数
	TokenMaxAge time.Duration
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
}

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &RedisConfig{}
var MongoSettings = &MongoConfig{}*/

type config struct {
	//自定义的配置
	Customize serverConfig
	//命令参数大于配置
	Flag     flagValue
	Server   initialize.ServerConfig
	Mail     initialize.MailConfig
	Database initialize.DatabaseConfig
	Redis    initialize.RedisConfig
	Log      initialize.LogConfig
	Consul   initialize.EtcdConfig
}

var Conf = &config{}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
		if c.Flag.Password != "" {
			c.Database.Password = c.Flag.Password
			c.Redis.Password = c.Database.Password
		}
		if c.Flag.MailPassword != "" {
			c.Mail.Password = c.Flag.MailPassword
		}
	}

	c.Customize.UploadMaxSize = c.Customize.UploadMaxSize * 1024 * 1024
	c.Customize.TokenMaxAge = time.Second * 60 * 60 * 24 * c.Customize.TokenMaxAge
}
