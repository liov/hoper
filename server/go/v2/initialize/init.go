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
	"sync"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/liov/hoper/go/v2/utils/configor"
	"github.com/liov/hoper/go/v2/utils/configor/nacos"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/fs/watch"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect"
	"github.com/liov/hoper/go/v2/utils/slices"
	"github.com/pelletier/go-toml"
)

//约定大于配置
var (
	InitConfig = &Init{}
	once sync.Once
)

func init() {
	flag.StringVar(&InitConfig.Env, "env", DEVELOPMENT, "环境")
	flag.StringVar(&InitConfig.ConfUrl, "conf", `./config.toml`, "配置文件路径")
	flag.StringVar(&InitConfig.AddConfigPath, "add", "", "额外配置文件名")
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

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
	InitKey     = "initialize"
)

type EnvConfig struct {
	NacosTenant string
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
	//附加配置，不对外公开的的配置,特定文件名,启用文件搜寻查找
	AddConfigPath string
	conf          NeedInit
	dao           Dao
	//closes     []interface{}
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
	fmt.Printf("Load config from: %s\n",InitConfig.ConfUrl)
	init.AddConfig()
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
			break
		}
	}
	init.ConfigCenter = &nacos.Config{
		Addr:   onceConfig.Nacos.Addr,
		Tenant: init.EnvConfig.NacosTenant,
		Group:  onceConfig.Nacos.Group,
		DataId: onceConfig.Module,
		Watch:  onceConfig.Nacos.Watch,
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

var alreadyRun struct {
	WatchConfig bool
	InitFunc    bool
}

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf Config, dao Dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}

	//逃逸到堆上了
	InitConfig.SetInit(conf, dao)
	InitConfig.LoadConfig()
	nacosClient := InitConfig.getConfig()
	go nacosClient.Listener(InitConfig.UnmarshalAndSet)
	log.Debug(conf)
	return func() {
		InitConfig.CloseDao()
		log.Sync()
	}
}

// 从nacos拉取配置并返回nacos client
func (init *Init) getConfig() *nacos.Client {
	nacosClient := init.ConfigCenter.NewClient()
	err := nacosClient.GetConfigAllInfoHandle(init.UnmarshalAndSet)
	if err != nil {
		log.Fatal(err)
	}
	return nacosClient
}

// 额外的配置文件
func (init *Init) AddConfig() {
	if init.AddConfigPath != "" {
		adCongPath, err := fs.FindFile(init.AddConfigPath)
		if err == nil {
			err := configor.New(&configor.Config{Debug: init.Env != PRODUCT}).
				Load(init.conf, adCongPath)
			if err != nil {
				log.Fatalf("配置错误: %v", err)
			}
		} else {
			log.Fatalf("找不到附加配置: %v", err)
		}
	}
}

//反射方法命名规范,P+优先级+方法名+(执行一次+Once)
func (init *Init) SetConfigAndDao() {
	confValue := reflect.ValueOf(init.conf).Elem()
	for i := 0; i < confValue.NumField(); i++ {
		if conf, ok := confValue.Field(i).Addr().Interface().(NeedInit); ok {
			conf.Custom()
		}
	}
	init.conf.Custom()
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
	if init.dao != nil {
		init.dao.Custom()
	}
}

func Watcher() {
	watcher, err := watch.New(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	watcher.Add(InitConfig.ConfUrl, fsnotify.Write, func() {
		InitConfig.CloseDao()
		InitConfig.SetConfigAndDao()
	})

	watcher.Add(".watch", fsnotify.Write, func() {
		watcher.Close()
	})
}

func (init *Init) Refresh() {
	InitConfig.CloseDao()
	InitConfig.SetConfigAndDao()
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
	init.Refresh()
}
