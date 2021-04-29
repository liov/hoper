package initialize

import (
	"flag"
	"fmt"
	"github.com/liov/hoper/go/v2/utils/slices"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/liov/hoper/go/v2/utils/configor"
	"github.com/liov/hoper/go/v2/utils/configor/nacos"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/pelletier/go-toml"
)

//约定大于配置
var (
	InitConfig = &Init{}
)

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
)

const (
	exprTag   = "expr"
	configTag = "config"
)

type EnvConfig struct {
	NacosTenant string
	//本地配置，特定文件名,启用文件搜寻查找
	LocalConfigName string
}

type BasicConfig struct {
	Module string
	NoInit []string
}

type Init struct {
	EnvConfig *EnvConfig
	// TODO: 接口抽象，可切换的配置中心
	ConfigCenter *nacos.Config
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

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf Config, dao Dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}

	//逃逸到堆上了
	InitConfig.SetInit(conf, dao)
	InitConfig.LoadConfig()
	return func() {
		InitConfig.CloseDao()
		log.Sync()
	}
}

func (init *Init) LoadConfig() *Init {
	onceConfig := struct {
		Dev, Test, Prod *EnvConfig
		Nacos           struct {
			Addr  string
			Group string
			Watch bool
		}
		BasicConfig
	}{}
	if _, err := os.Stat(InitConfig.ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置错误: 请确保可执行文件和配置目录在同一目录下")
	}
	err := configor.Load(&onceConfig, InitConfig.ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	fmt.Printf("Load config from: %s\n", InitConfig.ConfUrl)

	for i := range onceConfig.NoInit {
		onceConfig.NoInit[i] = strings.ToLower(onceConfig.NoInit[i])
	}
	init.BasicConfig = onceConfig.BasicConfig
	init.NoInit = onceConfig.NoInit

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
			if init.EnvConfig.NacosTenant != "" {
				init.ConfigCenter = &nacos.Config{
					Addr:   onceConfig.Nacos.Addr,
					Tenant: init.EnvConfig.NacosTenant,
					Group:  onceConfig.Nacos.Group,
					DataId: onceConfig.Module,
					Watch:  onceConfig.Nacos.Watch,
				}
				nacosClient := InitConfig.getNacosClient()
				go nacosClient.Listener(InitConfig.UnmarshalAndSet)
			} else if init.EnvConfig.LocalConfigName != "" {
				init.LocalConfig()
			} else {
				log.Fatal("没有发现配置")
			}
			break
		}
	}

	log.Debugf("Configuration:\n  %#v\n", init)
	return init
}

func (init *Init) SetInit(conf Config, dao Dao) {
	init.conf = conf
	init.dao = dao
}

type NeedInit interface {
	Custom()
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

// Custom
func (init *Init) inject() {
	var fieldNameDaoMap = make(map[string]interface{})
	setConfig(reflect.ValueOf(init.conf).Elem(), fieldNameDaoMap)
	init.conf.Custom()
	if init.dao == nil {
		return
	}
	setDao(reflect.ValueOf(init.dao).Elem(), fieldNameDaoMap)
	init.dao.Custom()
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
				conf.Custom()
			}
			if slices.StringContains(InitConfig.NoInit, strings.ToLower(typ.Field(i).Name)) {
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
			if daoValue.Type() == v.Field(i).Type() {
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
	InitConfig.CloseDao()
	InitConfig.inject()
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
