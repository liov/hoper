package initialize

import (
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/liov/hoper/go/v2/utils/configor"
	"github.com/liov/hoper/go/v2/utils/fs"
	"github.com/liov/hoper/go/v2/utils/fs/watch"
	"github.com/liov/hoper/go/v2/utils/log"
)

// 本地配置文件
func (init *Init) LocalConfig() {
	if init.EnvConfig.LocalConfigName != "" {
		adCongPath, err := fs.FindFile(init.EnvConfig.LocalConfigName)
		init.EnvConfig.LocalConfigName = adCongPath
		if err == nil {
			err := configor.New(&configor.Config{Debug: init.Env != PRODUCT}).
				Load(init.conf, adCongPath)
			if err != nil {
				log.Fatalf("配置错误: %v", err)
			}
			init.refresh()
			watcher()
		} else {
			log.Fatalf("找不到附加配置: %v", err)
		}
	}
}

func watcher() {
	watcher, err := watch.New(time.Second)
	if err != nil {
		log.Fatal(err)
	}
	watcher.Add(InitConfig.EnvConfig.LocalConfigName, fsnotify.Write, func() {
		err := configor.New(&configor.Config{Debug: InitConfig.Env == DEVELOPMENT}).
			Load(InitConfig.conf, InitConfig.EnvConfig.LocalConfigName)
		if err != nil {
			log.Fatalf("配置错误: %v", err)
		}
		InitConfig.refresh()
	})
}
