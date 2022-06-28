package local

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/utils/configor/local"
	"github.com/actliboy/hoper/server/go/lib/utils/fs"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
	"github.com/fsnotify/fsnotify"
)

type Local struct {
	local.Config
	ConfigName string
	ReloadType string `json:"reloadType" enum:"fsnotify,timer"` // 本地分为Watch和AutoReload，Watch采用系统调用通知，AutoReload定时器去查文件是否变更
}

// 本地配置
func (cc *Local) HandleConfig(handle func([]byte)) error {
	localConfigName := cc.ConfigName
	if localConfigName != "" {
		adCongPath, err := fs.FindFile(localConfigName)
		localConfigName = adCongPath
		if err == nil {
			var watch bool
			if cc.AutoReload && cc.ReloadType == "fsnotify" {
				cc.AutoReload = false
				watch = true
			}
			err := local.New(&cc.Config).
				Handle(handle, adCongPath)
			if err != nil {
				return fmt.Errorf("配置错误: %v", err)
			}
			if watch {
				go cc.watch(adCongPath, handle)
			}
		} else {
			return fmt.Errorf("找不到配置: %v", err)
		}

	}
	return nil
}

func (cc *Local) watch(adCongPath string, handle func([]byte)) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(adCongPath)
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			//log.Info("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				err = local.New(&cc.Config).
					Handle(handle, adCongPath)
				if err != nil {
					log.Errorf("配置错误: %v", err)
				}
				log.Info("modified file:", event.Name)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("error:", err)
		}
	}
}
