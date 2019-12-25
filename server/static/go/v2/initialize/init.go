package initialize

import (
	"flag"
	"os"
	"path"
	"reflect"
	"strings"
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
var AddConfig string

func init() {
	flag.StringVar(&ConfUrl, "c", "./config/config.toml", "配置文件夹路径")
	flag.StringVar(&AddConfig, "add", "add-config.toml", "额外配置文件名")
}

const (
	DEVELOPMENT = "dev"
	TEST        = "test"
	PRODUCT     = "prod"
)

//var closes = []interface{}{log.Sync}

type BasicConfig struct {
	Module string
	Env    string
	Volume fs.Dir
}

type Init struct {
	Module, Env  string
	HasAddConfig bool //附加配置，不对外公开的的配置,特定文件名,启用文件搜寻查找
	NoInit       []string
	conf         needInit
	dao          dao
	//closes     []interface{}
}

type Config interface {
	Generate(dao)
}

type needInit interface {
	Custom()
}

type config interface {
	//BasicConfig() *BasicConfig
	needInit
}

type dao interface {
	Close()
	needInit
}

var alreadyRun struct {
	WatchConfig bool
}

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf config, dao dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}
	init := &Init{conf: conf, dao: dao}
	init.config()
	init.apollo()
	init.setDao()
	if !alreadyRun.WatchConfig {
		//go Watcher(conf, dao)
		alreadyRun.WatchConfig = true
	}
	return func() {
		if dao != nil {
			dao.Close()
		}
		log.Sync()
		/*for _, f := range closes {
			res := reflect.ValueOf(f).Call(nil)
			if len(res) > 0 && res[0].IsValid() {
				log.Error(res[0].Interface())
			}
		}*/
	}
}

func (init *Init) config() {
	if _, err := os.Stat(ConfUrl); os.IsNotExist(err) {
		log.Fatalf("配置文件不存在：", err)
	}
	err := configor.New(&configor.Config{Debug: false}).
		Load(init, ConfUrl) //"./config/config.toml"
	dir, file := path.Split(ConfUrl)
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	err = configor.New(&configor.Config{Debug: init.Env != PRODUCT}).
		Load(init.conf, ConfUrl, dir+init.Env+path.Ext(file)) //"./config/{{env}}.toml"
	if err != nil {
		log.Fatalf("配置错误: %v", err)
	}
	if init.HasAddConfig {
		adCongPath, err := fs.FindFile(AddConfig)
		if err == nil {
			err := configor.New(&configor.Config{Debug: init.Env != PRODUCT}).
				Load(init.conf, adCongPath)
			if err != nil {
				log.Fatalf("配置错误: %v", err)
			}
		} else {
			log.Fatalf("配置错误: %v", err)
		}
	}
	init.conf.Custom()
}

func (init *Init) setDao() {
	value := reflect.ValueOf(init)
	noInit := strings.Join(init.NoInit, " ")
	typeOf := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		if strings.Contains(noInit, typeOf.Method(i).Name[2:]) {
			continue
		}
		if typeOf.Method(i).Type.NumOut() > 0 && init.dao == nil {
			continue
		}

		if res := value.Method(i).Call(nil); res != nil && len(res) > 0 {
			daoValue := reflect.ValueOf(init.dao).Elem()
			for j := range res {
				if res[j].IsValid() {
					reflect3.SetFieldValue(daoValue, res[j])
				}
			}
		}
	}
	if init.dao != nil {
		init.dao.Custom()
	}
}

func Watcher(conf config, dao dao) {
	watcher, err := watch.New(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	watcher.Add(ConfUrl, fsnotify.Write, func() {
		dao.Close()
		Start(conf, dao)
	})

	watcher.Add(".watch", fsnotify.Write, func() {
		watcher.Close()
	})
}
