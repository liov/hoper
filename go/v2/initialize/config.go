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
	flag.StringVar(&confUrl,"c", "./config", "配置文件夹路径")
}

func (i *Init) config() {
	err := configor.New(&configor.Config{Debug: false}).
		Load(i, confUrl+PathSeparator+"config"+ext) //"./config/config.toml"
	err = configor.New(&configor.Config{Debug: i.Env != PRODUCT}).
		Load(i.conf, confUrl+PathSeparator+"config"+ext,confUrl+PathSeparator+i.Env+ext) //"./config/{{env}}.toml"
	if err != nil {
		log.Errorf("配置错误: %v", err)
		os.Exit(10)
	}
	i.conf.Custom()
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

type ElasticConfig struct {
	Host string
}

type NsqConfig struct {
	Host string
	model int8 //0生产者，1消费者，2所有
	topic string
	channel string
}

type MongoConfig struct {
	URL      string
	Database string
}

type LogConfig struct {
	Level    int8
	SaveDir fs.Dir
	OutputPaths []string //日志文件路径
	ErrOutputPaths []string
}

type MailConfig struct {
	Host     string
	Port     string
	User     string
	Password string
}