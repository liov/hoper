package initialize

import (
	"flag"
	"os"
	"time"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/log"
)

const (
	ext           = ".toml"
	PathSeparator = "/"
)

var confUrl string

func init() {
	flag.StringVar(&confUrl, "c", "./config", "配置文件夹路径")
}

func (i *Init) config() {
	err := configor.New(&configor.Config{Debug: false}).
		Load(i, confUrl+PathSeparator+"config"+ext) //"./config/config.toml"
	err = configor.New(&configor.Config{Debug: i.Env != PRODUCT}).
		Load(i.conf, confUrl+PathSeparator+"config"+ext, confUrl+PathSeparator+i.Env+ext) //"./config/{{env}}.toml"
	if err != nil {
		log.Errorf("配置错误: %v", err)
		os.Exit(10)
	}
	i.conf.Custom()
}

type BasicConfig struct {
	Module     string
	Env        string
	Volume     fs.Dir
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
	Port         int32
}

type RedisConfig struct {
	Addr        string
	Password    string
	MaxIdle     int
	MaxActive   int
	IdleTimeout time.Duration
}

type ElasticConfig struct {
	Host string
	Port int32
}

type NsqConfig struct {
	Addr    string
	Model   int8 //0生产者，1消费者，2所有
	Topic   string
	Channel string
}

type LogConfig struct {
	Level       int8
	Skip        bool
	OutputPaths map[string][]string
}

type MailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}

type KafkaConfig struct {
	Model    int8 //0生产者，1消费者，2所有
	Topic    string
	ProdAddr []string
	ConsAddr []string
}
