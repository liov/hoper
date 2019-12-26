package config

import (
	"runtime"
	"time"

	"github.com/liov/hoper/go/v2/initialize"
	"github.com/liov/hoper/go/v2/utils/fs"
)

type serverConfig struct {
	Protocol     string
	Domain       string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PassSalt    string
	TokenMaxAge int64
	TokenSecret string
	JwtSecret   string
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
	//必须
	initialize.BasicConfig
	//自定义的配置
	Server serverConfig
	//命令参数大于配置
	Flag     flagValue
	Mail     initialize.MailConfig
	Database initialize.DatabaseConfig
	Redis    initialize.RedisConfig
	Log      initialize.LogConfig
	Apollo   initialize.ApolloConfig
	Consul   initialize.ConsulConfig
}

var Conf = &config{}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}

func (c *config) Custom() {
	if runtime.GOOS == "windows" {
		c.Server.LuosimaoAPIKey = ""
		if c.Flag.Password != "" {
			c.Database.Password = c.Flag.Password
			c.Redis.Password = c.Database.Password
		}
		if c.Flag.MailPassword != "" {
			c.Mail.Password = c.Flag.MailPassword
		}
	}

	c.Server.UploadMaxSize = c.Server.UploadMaxSize * 1024 * 1024
	c.Server.ReadTimeout = c.Server.ReadTimeout * time.Second
	c.Server.WriteTimeout = c.Server.WriteTimeout * time.Second
}
