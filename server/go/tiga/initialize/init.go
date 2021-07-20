package initialize

import (
	"flag"
	"fmt"
	"github.com/liov/hoper/v2/tiga/initialize/conf_center"
	"github.com/liov/hoper/v2/utils/slices"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/liov/hoper/v2/utils/configor"
	"github.com/liov/hoper/v2/utils/log"
	"github.com/pelletier/go-toml"
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

const (
	exprTag   = "expr"
	configTag = "config"
)

type EnvConfig struct {
	conf_center.ConfigCenterEnvConfig
}

type BasicConfig struct {
	Module   string
	NoInject []string
}

type Init struct {
	EnvConfig *EnvConfig
	BasicConfig
	Env, ConfUrl string
	conf         NeedInit
	dao          Dao
	//closes     []interface{}
}

func init() {
	flag.StringVar(&InitConfig.Env, "env", DEVELOPMENT, "环境")
	flag.StringVar(&InitConfig.ConfUrl, "conf", `./config.toml`, "配置文件路径")
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

func Start(conf Config, dao Dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}

	//逃逸到堆上了
	init := NewInit(conf, dao)
	init.LoadConfig()
	return func() {
		init.CloseDao()
		log.Sync()
	}
}

func (init *Init) LoadConfig() *Init {
	onceConfig := struct {
		Dev, Test, Prod *EnvConfig
		ConfigCenter    conf_center.ConfigCenterConfig
		BasicConfig
	}{}
	if _, err := os.Stat(init.ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置错误: 请确保可执行文件和配置目录在同一目录下")
	}
	err := configor.Load(&onceConfig, init.ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	fmt.Printf("Load config from: %s\n", init.ConfUrl)

	for i := range onceConfig.NoInject {
		onceConfig.NoInject[i] = strings.ToLower(onceConfig.NoInject[i])
	}
	init.BasicConfig = onceConfig.BasicConfig
	init.NoInject = onceConfig.NoInject

	value := reflect.ValueOf(&onceConfig).Elem()
	typ := reflect.TypeOf(&onceConfig).Elem()

	tmpTyp := reflect.TypeOf(&EnvConfig{})
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == tmpTyp && strings.ToUpper(typ.Field(i).Name) == strings.ToUpper(init.Env) {
			/*tmpConfig = value.Field(i).Interface().(*nacos.Config)
			//真·深度复制
			data,_:=json.Marshal(tmpConfig)
			if err:=json.Unmarshal(data,init.EnvConfig);err!=nil{
				log.Fatal(err)
			}*/
			//会被回收,也可能是被移动了？
			init.EnvConfig = &(*value.Field(i).Interface().(*EnvConfig))
			onceConfig.ConfigCenter.ConfigCenter(init.EnvConfig.ConfigCenterEnvConfig, init.Module, init.Env != PRODUCT).HandleConfig(init.UnmarshalAndSet)
			break
		}
	}

	log.Debugf("Configuration:\n  %#v\n", init)
	return init
}

func NewInit(conf Config, dao Dao) *Init {
	init := &Init{
		Env: InitConfig.Env, ConfUrl: InitConfig.ConfUrl,
		conf: conf, dao: dao}
	InitConfig = init
	return init
}

func (init *Init) SetInit(conf Config, dao Dao) {
	init.conf = conf
	init.dao = dao
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

type Generate interface {
	Generate() interface{}
}

// Customize
func (init *Init) inject() {
	var fieldNameDaoMap = make(map[string]interface{})
	setConfig(reflect.ValueOf(init.conf).Elem(), fieldNameDaoMap)
	init.conf.Init()
	if init.dao == nil {
		return
	}
	setDao(reflect.ValueOf(init.dao).Elem(), fieldNameDaoMap)
	init.dao.Init()
}

func setConfig(v reflect.Value, fieldNameDaoMap map[string]interface{}) {
	if !v.IsValid() {
		return
	}
	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			setConfig(v.Field(i).Elem(), fieldNameDaoMap)
		case reflect.Struct:
			setConfig(v.Field(i), fieldNameDaoMap)
		}
		if v.Field(i).Addr().CanInterface() {
			inter := v.Field(i).Addr().Interface()
			if conf, ok := inter.(NeedInit); ok {
				conf.Init()
			}
			if slices.StringContains(InitConfig.NoInject, strings.ToLower(typ.Field(i).Name)) {
				continue
			}
			if conf, ok := inter.(Generate); ok {
				ret := conf.Generate()
				fieldNameDaoMap[strings.ToLower(typ.Field(i).Name)] = ret
			}
		}
	}
}

func setDao(v reflect.Value, fieldNameDaoMap map[string]interface{}) {

	if !v.IsValid() {
		return
	}
	typ := v.Type()

	for i := 0; i < v.NumField(); i++ {
		var dao interface{}
		var ok bool
		if confName, cok := typ.Field(i).Tag.Lookup("config"); cok {
			dao, ok = fieldNameDaoMap[strings.ToLower(confName)]
		}
		if !ok {
			dao, ok = fieldNameDaoMap[strings.ToLower(typ.Field(i).Name)]
		}
		if ok {
			daoValue := reflect.ValueOf(dao)
			if daoValue.Type().AssignableTo(v.Field(i).Type()) || daoValue.Type().Implements(v.Field(i).Type()) {
				v.Field(i).Set(daoValue)
			}
			continue
		}
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			setDao(v.Field(i).Elem(), fieldNameDaoMap)
		case reflect.Struct:
			setDao(v.Field(i), fieldNameDaoMap)
		}
	}
}
func (init *Init) refresh() {
	init.CloseDao()
	init.inject()
}

func (init *Init) CloseDao() {
	if init.dao != nil {
		init.dao.Close()
	}
}

func (init *Init) UnmarshalAndSet(bytes []byte) {
	toml.Unmarshal(bytes, init.conf)
	init.refresh()
}
