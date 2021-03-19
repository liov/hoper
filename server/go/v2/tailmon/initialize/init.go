package initialize

import (
	"flag"
	"fmt"
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
	"github.com/liov/hoper/go/v2/utils/reflect"
	"github.com/liov/hoper/go/v2/utils/slices"
	"github.com/pelletier/go-toml"
)

//约定大于配置
var (
	InitConfig = &Init{}
	alreadyRun struct {
		WatchConfig bool
		InitFunc    bool
	}
)

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
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



// Custom
func (init *Init) setConfig() {
	setConfig(reflect.ValueOf(init.conf).Elem())
	init.conf.Custom()
}

func setConfig(v reflect.Value) {
	if !v.IsValid() {
		return
	}
	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).Addr().CanInterface() {
			if conf, ok := v.Field(i).Addr().Interface().(NeedInit); ok {
				conf.Custom()
			}
		}
		switch v.Field(i).Kind() {
		case reflect.Ptr:
			setConfig(v.Field(i).Elem())
		case reflect.Struct:
			setConfig(v.Field(i))
		}
	}
}

//反射方法命名规范,P+优先级+方法名+(执行一次+Once)
func (init *Init) setDao() {
	if init.dao == nil {
		return
	}
	value := reflect.ValueOf(init)
	typeOf := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		methodName := typeOf.Method(i).Name
		if strings.Contains(methodName, "Once") {
			if alreadyRun.InitFunc || slices.StringContains(init.NoInit, methodName[2:len(methodName)-4]) {
				continue
			}
		}
		if !strings.HasPrefix(methodName, "P") || slices.StringContains(init.NoInit, methodName[2:]) {
			continue
		}

		if res := value.Method(i).Call(nil); len(res) > 0 {
			daoValue := reflect.ValueOf(init.dao).Elem()
			for j := range res {
				if res[j].IsValid() {
					reflecti.SetFieldValue(daoValue, res[j])
				}
			}
		}
	}
	alreadyRun.InitFunc = true
	init.dao.Custom()
}

func (init *Init) refresh() {
	InitConfig.CloseDao()
	InitConfig.setConfig()
	InitConfig.setDao()
}

func (init *Init) Unmarshal(bytes []byte) {
	toml.Unmarshal(bytes, init.conf)
}

func (init *Init) CloseDao() {
	if init.dao != nil {
		init.dao.Close()
	}
}

func (init *Init) UnmarshalAndSet(bytes []byte) {
	log.Debug(string(bytes))
	init.Unmarshal(bytes)
	init.refresh()
}
