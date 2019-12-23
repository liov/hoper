package initialize

import (
	"flag"
	"os"
	"path"
	"time"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/log"
)

var ConfUrl string

func init() {
	flag.StringVar(&ConfUrl, "c", "./config/config.toml", "配置文件夹路径")
}

func (i *Init) config() {
	if _, err := os.Stat(ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在：", err)
	}
	err := configor.New(&configor.Config{Debug: false}).
		Load(i, ConfUrl) //"./config/config.toml"
	dir, file := path.Split(ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	err = configor.New(&configor.Config{Debug: i.Env != PRODUCT}).
		Load(i.conf, ConfUrl, dir+i.Env+path.Ext(file)) //"./config/{{env}}.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	if i.Additional != "" {
		err := configor.New(&configor.Config{Debug: i.Env != PRODUCT}).
			Load(i.conf, i.Additional)
		if err != nil {
			log.Fatalf("配置错误: %v", err)
		}
	}
	i.conf.Custom()
}

type BasicConfig struct {
	Module string
	Env    string
	Volume fs.Dir
}

type DatabaseConfig struct {
	Type         string
	User         string
	Password     string
	Host         string
	Charset      string
	Database     string
	TimeFormat   string
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
	From     string
	Password string
}

type KafkaConfig struct {
	Model    int8 //0生产者，1消费者，2所有
	Topic    string
	ProdAddr []string
	ConsAddr []string
}

type ApolloConfig struct {
	Addr      string
	AppId     string `json:"appId"`
	Cluster   string `json:"cluster"`
	IP        string `json:"ip"`
	NameSpace []string
}
