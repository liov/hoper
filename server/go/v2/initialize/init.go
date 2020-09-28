package initialize

import (
	"flag"
	"net"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/nacos"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/liov/hoper/go/v2/utils/watch"
)

//约定大于配置
var ConfUrl = "./config.toml"

var Env string

//附加配置，不对外公开的的配置,特定文件名,启用文件搜寻查找
var AddConfig string

var InitConfig *Init

var once sync.Once

func init() {
	flag.StringVar(&Env, "env", DEVELOPMENT, "环境")
	flag.StringVar(&AddConfig, "add", "", "额外配置文件名")
	agent := flag.Bool("agent", false, "是否启用代理")
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
	Tenant string
}

type BasicConfig struct {
	Module string
	Env    string
}

type BasicNacosConfig struct {
	NacosAddr  string
	NacosGroup string
	//NacosDataId string DataId与Moudle相同
}

type OnceConfig struct {
	Dev, Test, Prod *EnvConfig
	BasicNacosConfig
	BasicConfig
	NoInit []string
}

type Init struct {
	EnvConfig   *EnvConfig
	NacosConfig *nacos.Config
	BasicConfig
	NoInit []string
	conf   NeedInit
	dao    Dao
	//closes     []interface{}
}

func NewInitWithLoadConfig(conf Config, dao Dao) *Init {
	onceConfig := OnceConfig{}
	if _, err := os.Stat(ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置错误: 请确保可执行文件和配置目录在同一目录下")
	}
	err := configor.New(&configor.Config{Debug: true}).
		Load(&onceConfig, ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	init := NewInit(conf, dao)

	init.BasicConfig = onceConfig.BasicConfig
	init.Env = Env
	init.NoInit = onceConfig.NoInit

	value := reflect.ValueOf(&onceConfig).Elem()
	typ := reflect.TypeOf(&onceConfig).Elem()
	var tmpConfig = EnvConfig{}
	tmpTyp := reflect.TypeOf(&tmpConfig)
	for i := 0; i < typ.NumField(); i++ {
		if typ.Field(i).Type == tmpTyp && strings.ToUpper(typ.Field(i).Name) == strings.ToUpper(init.Env) {
			/*tmpConfig = value.Field(i).Interface().(*nacos.Config)
			//真·深度复制
			data,_:=json.Marshal(tmpConfig)
			if err:=json.Unmarshal(data,init.EnvConfig);err!=nil{
				log.Fatal(err)
			}*/
			//会被回收,也可能是被移动了？
			tmpConfig = *value.Field(i).Interface().(*EnvConfig)
			init.EnvConfig = &tmpConfig
			break
		}
	}
	init.NacosConfig = &nacos.Config{
		Addr:   onceConfig.NacosAddr,
		Tenant: init.EnvConfig.Tenant,
		Group:  onceConfig.NacosGroup,
		DataId: onceConfig.Module,
	}
	return init
}

func NewInit(conf Config, dao Dao) *Init {
	return &Init{conf: conf, dao: dao}
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
	init := NewInitWithLoadConfig(conf, dao)
	InitConfig = init
	nacosClient := init.getConfig()
	go nacosClient.Listener(init.UnmarshalAndSet)

	return func() {
		InitConfig.CloseDao()
		log.Sync()
	}
}

func (init *Init) getConfig() *nacos.Client {
	nacosClient := init.NacosConfig.NewClient()
	err := nacosClient.GetConfigAllInfoHandle(init.UnmarshalAndSet)
	if err != nil {
		log.Fatal(err)
	}
	return nacosClient
}

func (init *Init) AddConfig(addConfig string) *Init {
	if addConfig != "" {
		adCongPath, err := fs.FindFile(addConfig)
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
	return init
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
	noInit := strings.Join(init.NoInit, " ")
	typeOf := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		methodName := typeOf.Method(i).Name
		if strings.Contains(typeOf.Method(i).Name, "Once") {
			if strings.Contains(noInit, methodName[2:len(methodName)-4]) || alreadyRun.InitFunc {
				continue
			}
		}
		if strings.Contains(noInit, methodName[2:]) || !strings.HasPrefix(methodName, "P") {
			continue
		}

		if res := value.Method(i).Call(nil); len(res) > 0 {
			daoValue := reflect.ValueOf(init.dao).Elem()
			for j := range res {
				if res[j].IsValid() {
					reflect3.SetFieldValue(daoValue, res[j])
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
	watcher.Add(ConfUrl, fsnotify.Write, func() {
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
	init.Unmarshal(bytes)
	init.Refresh()
}
