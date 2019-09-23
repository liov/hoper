package initialize

import (
	"flag"
	"os"
	"reflect"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/utls/log"
)

func (i *Init) config(config interface{}) {
	confUrl := flag.String("c", "./config/config.toml", "配置文件路径")
	err := configor.New(&configor.Config{Debug: false}).Load(i, *confUrl)
	err = configor.New(&configor.Config{Debug: i.Env != PRODUCT}).Load(config, *confUrl)
	if err != nil {
		log.Errorf("配置错误: %v", err)
		os.Exit(10)
	}
	value := reflect.ValueOf(config)
	value.MethodByName("Set").Call([]reflect.Value{reflect.ValueOf(i.Env)})
}
