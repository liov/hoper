package config

import (
	"flag"
	"runtime"
	"time"

	"go.uber.org/zap/zapcore"
)

type ServerConfig struct {
	Env          string
	HttpPort     string
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

	MailHost     string
	MailPort     string
	MailUser     string
	MailPassword string

	LuosimaoVerifyURL string
	LuosimaoAPIKey    string

	QrCodeSavePath string //二维码保存路径
	PrefixUrl      string
	FontSavePath   string //字体保存路径

	CrawlerName string //爬虫
}

type DatabaseConfig struct {
	Type         string
	User         string
	Password     string
	Host         string
	Charset      string
	Database     string
	TablePrefix  string
	MaxIdleConns int
	MaxOpenConns int
	Port         int
}

type RedisConfig struct {
	Host        string
	Port        int
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type MongoConfig struct {
	URL      string
	Database string
}

type LogConfig struct {
	Level    zapcore.Level
	FilePath []string //日志文件路径
}

/*var ServerSettings = &ServerConfig{}
var DatabaseSettings = &DatabaseConfig{}
var RedisSettings = &RedisConfig{}
var MongoSettings = &MongoConfig{}*/

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	Mongo    MongoConfig
	Log      LogConfig
}

var Conf = &Config {}

type duration struct {
	time.Duration
}

func (d *duration) UnmarshalText(text []byte) error {
	var err error
	d.Duration, err = time.ParseDuration(string(text))
	return err
}
//固定函数名，init时反射用
func (c *Config) Set(env string)  {
	if runtime.GOOS == "windows" {
		c.Server.LuosimaoAPIKey = ""
		c.Redis.Password = ""
	} else {
		flag.StringVar(&c.Database.Password, "p", c.Database.Password, "password")
		flag.StringVar(&c.Server.MailPassword, "mp", c.Server.MailPassword, "password")
		flag.Parse()
		c.Redis.Password = c.Database.Password
		c.Server.Env = env
	}

	c.Server.UploadMaxSize = c.Server.UploadMaxSize * 1024 * 1024
	c.Server.ReadTimeout = c.Server.ReadTimeout * time.Second
	c.Server.WriteTimeout = c.Server.WriteTimeout * time.Second
	c.Redis.IdleTimeout = c.Redis.IdleTimeout * time.Second
}