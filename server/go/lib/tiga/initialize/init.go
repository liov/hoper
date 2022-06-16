package initialize

import (
	"flag"
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/conf_center"
	"github.com/actliboy/hoper/server/go/lib/utils/configor/local"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

//约定大于配置
var (
	InitConfig = &Init{}
)

type Env string

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
	InitKey     = "initialize"
)

type ConfigCenterConfig struct {
	NoInject      []string
	InjectVersion int8
	conf_center.ConfigCenterConfig
}

type BasicConfig struct {
	Module string
}

type FileConfig struct {
	BasicConfig
	Dev, Test, Prod *ConfigCenterConfig
}

//TODO: 优先级高于EnvConfig
type FlagConfig struct {
	Env, ConfUrl string
}

//TODO: 优先级高于FileConfig
type EnvConfig FlagConfig

type Init struct {
	Env, ConfUrl string
	BasicConfig
	ConfigCenterConfig *ConfigCenterConfig
	confM              map[string]interface{}
	conf               NeedInit
	dao                Dao
	//closes     []interface{}
	deferf []func()
}

func flaginit() {
	if flag.Parsed() {
		return
	}
	flag.StringVar(&InitConfig.Env, "env", DEVELOPMENT, "环境")

	InitConfig.ConfUrl = "./config.toml"
	if _, err := os.Stat(InitConfig.ConfUrl); os.IsNotExist(err) {
		InitConfig.ConfUrl = "./config/config.toml"
	}
	flag.StringVar(&InitConfig.ConfUrl, "conf", InitConfig.ConfUrl, "配置文件路径,默认./config.toml或./config/config.toml")

	agent := flag.Bool("agent", false, "是否启用代理")
	testing.Init()
	flag.Parse()
	if *agent {
		proxyURL, _ := url.Parse("socks5://localhost:1080")
		http.DefaultClient.Transport = &http.Transport{
			Proxy: http.ProxyURL(proxyURL),
			DialContext: (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext,
		}
	}
}

func Start(conf Config, dao Dao, notinit ...string) func(deferf ...func()) {
	if conf == nil {
		log.Fatalf("配置不能为空")
	}
	flaginit()
	//逃逸到堆上了
	init := NewInit(conf, dao)
	init.LoadConfig(notinit...)
	return func(deferf ...func()) {
		init.CloseDao()
		log.Sync()
		for _, f := range deferf {
			f()
		}
		for _, f := range init.deferf {
			f()
		}
	}
}

func (init *Init) LoadConfig(notinit ...string) *Init {
	onceConfig := FileConfig{}
	if _, err := os.Stat(init.ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置错误: 请确保可执行文件和配置文件在同一目录下或在config目录下或指定配置文件")
	}
	err := local.Load(&onceConfig, init.ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	fmt.Printf("Load config from: %s\n", init.ConfUrl)

	init.BasicConfig = onceConfig.BasicConfig

	value := reflect.ValueOf(&onceConfig).Elem()
	typ := reflect.TypeOf(&onceConfig).Elem()

	tmpTyp := reflect.TypeOf(&ConfigCenterConfig{})
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == tmpTyp && strings.ToUpper(typ.Field(i).Name) == strings.ToUpper(init.Env) {
			/*tmpConfig = value.Field(i).Interface().(*nacos.Config)
			//真·深度复制
			data,_:=json.Marshal(tmpConfig)
			if err:=json.Unmarshal(data,init.ConfigCenterConfig);err!=nil{
				log.Fatal(err)
			}*/
			//会被回收,也可能是被移动了？
			init.ConfigCenterConfig = &(*value.Field(i).Interface().(*ConfigCenterConfig))
			break
		}
	}

	for i := range init.ConfigCenterConfig.NoInject {
		init.ConfigCenterConfig.NoInject[i] = strings.ToUpper(init.ConfigCenterConfig.NoInject[i])
	}
	for i := range notinit {
		init.ConfigCenterConfig.NoInject = append(init.ConfigCenterConfig.NoInject, strings.ToUpper(notinit[i]))
	}
	if init.ConfigCenterConfig.InjectVersion == 1 {
		init.ConfigCenterConfig.ConfigCenter(init.Module, init.Env != PRODUCT).HandleConfig(init.UnmarshalAndSet)
	} else {
		init.ConfigCenterConfig.ConfigCenter(init.Module, init.Env != PRODUCT).HandleConfig(init.UnmarshalAndSetV2)
	}

	log.Debugf("Configuration:  %#v", init.conf)
	return init
}

func NewInit(conf Config, dao Dao) *Init {
	init := &Init{
		Env: InitConfig.Env, ConfUrl: InitConfig.ConfUrl,
		confM: map[string]interface{}{},
		conf:  conf, dao: dao}
	InitConfig = init
	return init
}

func (init *Init) SetInit(conf Config, dao Dao) {
	init.conf = conf
	init.dao = dao
}

func (init *Init) RegisterDeferFunc(deferf ...func()) {
	init.deferf = append(init.deferf, deferf...)
}

type NeedInit interface {
	Init()
}

type Config interface {
	NeedInit
}

type Dao interface {
	Close()
	NeedInit
}

type DaoField interface {
	Config() any
	SetEntity(any)
}

type Generate interface {
	Generate() any
}

func (init *Init) CloseDao() {
	if init.dao != nil {
		init.dao.Close()
	}
}

func (init *Init) Config() Config {
	return init.conf
}
