package initialize

import (
	"flag"
	"path"
	"path/filepath"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/liov/hoper/go/v2/utils/watch"
)

type EnvVar map[string]string

var ConfUrl string

//附加配置，不对外公开的的配置,特定文件名,启用文件搜寻查找
var AddConfig string

func init() {
	flag.StringVar(&ConfUrl, "c", "./Config/Config.toml", "配置文件夹路径")
	flag.StringVar(&AddConfig, "add", "", "额外配置文件名")
}

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
	InitKey     = "initialize"
)

//var closes = []interface{}{log.Sync}

type BasicConfig struct {
	Module string
	Env    string
	Volume fs.Dir
}

type Init struct {
	Module, Env string
	NoInit      []string
	conf        needInit
	dao         Dao
	//closes     []interface{}
}

func NewInit(conf Config, dao Dao) *Init {
	init := &Init{conf: conf, dao: dao}
	err := configor.New(&configor.Config{Debug: false}).
		Load(init, ConfUrl) //"./Config/Config.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	return init
}

//配置作用于dao，但是这么写dao不直观，无法跳转，不利于阅读
type config interface {
	Generate(Dao)
}

type needInit interface {
	Custom()
}

type Config interface {
	needInit
}

type Dao interface {
	Close()
	needInit
}

var alreadyRun struct {
	WatchConfig bool
	InitFunc    bool
}

var once sync.Once

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf Config, dao Dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}
	init := NewInit(conf, dao)
	init.config()
	//从config到dao的过渡
	init.setDao()
	//go Watcher(conf, Dao)
	return func() {
		if dao != nil {
			dao.Close()
		}
		log.Sync()
		/*for _, f := range closes {
			res := reflect.ValueOf(f).Caller(nil)
			if len(res) > 0 && res[0].IsValid() {
				log.Error(res[0].Interface())
			}
		}*/
	}
}

func (init *Init) config() {
	dir, file := filepath.Split(ConfUrl)
	err := configor.New(&configor.Config{Debug: init.Env != PRODUCT}).
		Load(init.conf, ConfUrl, dir+init.Env+path.Ext(file)) //"./Config/{{env}}.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	if AddConfig != "" {
		adCongPath, err := fs.FindFile(AddConfig)
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
func (init *Init) setDao() {
	init.conf.Custom()
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
		if strings.Contains(noInit, methodName[2:]) {
			continue
		}
		if typeOf.Method(i).Type.NumOut() > 0 && init.dao == nil {
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

func Watcher(conf Config, dao Dao) {
	watcher, err := watch.New(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	watcher.Add(ConfUrl, fsnotify.Write, func() {
		dao.Close()
		init := NewInit(conf, dao)
		init.config()
		init.setDao()
	})

	watcher.Add(".watch", fsnotify.Write, func() {
		watcher.Close()
	})
}

func Refresh(conf Config, dao Dao) {
	init := NewInit(conf, dao)
	dao.Close()
	init.setDao()
}
