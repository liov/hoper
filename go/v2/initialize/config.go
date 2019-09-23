package initialize

import (
	"flag"
	"os"
	"reflect"
	"time"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/utils/log"
	"go.uber.org/zap/zapcore"
)

const (
	ext           = ".toml"
	PathSeparator = "/"
)

func (i *Init) config(config interface{}) {
	confUrl := flag.String("c", "./config", "配置文件夹路径")
	flag.Parse()
	err := configor.New(&configor.Config{Debug: false}).
		Load(i, *confUrl+PathSeparator+"config"+ext) //"./config/config.toml"
	err = configor.New(&configor.Config{Debug: i.Env != PRODUCT}).
		Load(config, *confUrl+PathSeparator+i.Env+ext) //"./config/{{env}}.toml"
	if err != nil {
		log.Errorf("配置错误: %v", err)
		os.Exit(10)
	}
	value := reflect.ValueOf(config)
	value.MethodByName("Set").Call([]reflect.Value{reflect.ValueOf(i.Env)})
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
