package initialize

import (
	"flag"
	"os"
	"path"

	"github.com/jinzhu/configor"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/log"
)

var ConfUrl string

func init() {
	flag.StringVar(&ConfUrl, "c", "./config/config.toml", "配置文件夹路径")
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
