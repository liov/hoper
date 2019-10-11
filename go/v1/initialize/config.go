package initialize

import (
	"flag"
	"os"
	"time"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v1/utils/log"
	"go.uber.org/zap/zapcore"
)

const (
	ext           = ".toml"
	PathSeparator = "/"
)

var confUrl string

func init() {
	flag.StringVar(&confUrl,"c", "./config", "配置文件夹路径")
}

type Config interface {
	Custom()
}

func (i *Init) config(config Config) {
	if flag.Parsed(){
		flag.Parse()
	}
	err := configor.New(&configor.Config{Debug: false}).
		Load(i, confUrl+PathSeparator+"config"+ext) //"./config/config.toml"
	err = configor.New(&configor.Config{Debug: i.Env != PRODUCT}).
		Load(config, confUrl+PathSeparator+"config"+ext,confUrl+PathSeparator+i.Env+ext) //"./config/{{env}}.toml"
	if err != nil {
		log.Errorf("配置错误: %v", err)
		os.Exit(10)
	}
	config.Custom()
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
