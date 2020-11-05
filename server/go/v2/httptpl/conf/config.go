package conf

import (
	"runtime"
	"time"

	"github.com/liov/hoper/go/v2/initialize"
)

type serverConfig struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PassSalt        string
	TokenMaxAge     int64
	TokenSecret     string
	JwtSecret       string
	PageSize        int8
	RuntimeRootPath string

	UploadDir      string
	UploadPath     string
	UploadMaxSize  int
	UploadAllowExt []string

	LogSavePath string
	LogSaveName string
	LogFileExt  string
	TimeFormat  string

	SiteName string
	Host     string

	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSavePath string //二维码保存路径
	PrefixUrl      string
	FontSavePath   string //字体保存路径

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
	Mail     initialize.MailConfig
	Server   initialize.ServerConfig
	Database initialize.DatabaseConfig
	Redis    initialize.RedisConfig
	Log      initialize.LogConfig
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

//固定函数名，init时反射用
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
	c.Server.ReadTimeout = c.Server.ReadTimeout * time.Second
	c.Server.WriteTimeout = c.Server.WriteTimeout * time.Second
	c.Redis.IdleTimeout = c.Redis.IdleTimeout * time.Second
}
