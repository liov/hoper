package initialize

import (
	"flag"
	"os"
	"path"
	"reflect"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

var ConfUrl string

func init() {
	flag.StringVar(&ConfUrl, "c", "./config/config.toml", "配置文件夹路径")
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
	BasicConfig
	NoInit        []string
	HasAdditional bool //附加配置，不对外公开的的配置,特定文件名,启用文件搜寻查找
	conf          needInit
	dao           dao
	//closes     []interface{}
}

type Config interface {
	Generate(dao)
}

type needInit interface {
	Custom()
}

type dao interface {
	Close()
	needInit
}

//init函数命名规则，P+数字（优先级）+ 功能名
func Start(conf needInit, dao dao) func() {
	if !flag.Parsed() {
		flag.Parse()
	}
	init := &Init{conf: conf, dao: dao}
	init.config()
	value := reflect.ValueOf(init)
	noInit := strings.Join(init.NoInit, " ")
	typeOf := value.Type()
	for i := 0; i < value.NumMethod(); i++ {
		if strings.Contains(noInit, typeOf.Method(i).Name[2:]) {
			continue
		}
		if typeOf.Method(i).Type.NumOut() > 0 && dao == nil {
			continue
		}

		if res := value.Method(i).Call(nil); res != nil && len(res) > 0 {
			daoValue := reflect.ValueOf(dao).Elem()
			for j := range res {
				if res[j].IsValid() {
					reflect3.SetFieldValue(daoValue, res[j])
				}
			}
		}
	}
	if dao != nil {
		dao.Custom()
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
	if i.HasAdditional {
		adCongPath, err := fs.FindFile2("add-config.toml", 5, 1)
		if err == nil {
			err := configor.New(&configor.Config{Debug: i.Env != PRODUCT}).
				Load(i.conf, adCongPath[0])
			if err != nil {
				log.Fatalf("配置错误: %v", err)
			}
		} else {
			log.Fatalf("配置错误: %v", err)
		}
	}
	i.conf.Custom()
}
