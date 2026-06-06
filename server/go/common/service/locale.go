package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/hopeio/gox/log"
	pathx "github.com/hopeio/gox/os/fs/path"
	"github.com/liov/hoper/server/go/global"
	"github.com/liov/hoper/server/go/protobuf/common"
	"go.uber.org/zap"
)

var localeMessages = map[string]map[string]string{"zh-CN": {}, "en": {}}
var lock sync.RWMutex
var editTimes int

func (*CommonService) Locale(ctx context.Context, req *common.LocaleReq) (*common.LocaleResp, error) {
	loadLocaleConfigs()
	locale := req.Locale
	if locale == "zh-CN" {
		locale = "zh-Hans"
	}
	lock.RLock()
	defer lock.RUnlock()
	messages := localeMessages[locale]
	if len(messages) == 0 {
		locale = global.Conf.Locale.Default
		messages = localeMessages[locale]
	}
	return &common.LocaleResp{Locale: locale, Messages: messages}, nil
}

func loadLocaleConfigs() {
	lock.Lock()
	defer lock.Unlock()
	editTimes++
	if editTimes > 1 {
		return
	}
	dir := global.Conf.Locale.Dir
	files, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, file := range files {
		messages := loadLocaleFile(filepath.Join(dir, file.Name()))
		if len(messages) == 0 {
			continue
		}
		localeMessages[pathx.FileNoExt(file.Name())] = messages
	}
	watch, err := fsnotify.NewWatcher()
	if err != nil {
		log.Errorw("fsnotify.NewWatcher error", zap.Error(err), zap.String("dir", dir))
		return
	}
	absDir, err := filepath.Abs(dir)
	if err != nil {
		log.Errorw("filepath.Abs error", zap.Error(err), zap.String("dir", dir))
		return
	}
	err = watch.Add(absDir)
	if err != nil {
		log.Errorw("watch.Add error", zap.Error(err), zap.String("dir", dir))
		return
	}
	go func() {
		defer watch.Close()
		defer lock.Unlock()
		for {
			select {
			case event, ok := <-watch.Events:
				log.Infow("locale file event", zap.Any("event", event))
				if !ok {
					return
				}
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					lock.Lock()
					messages := loadLocaleFile(filepath.Join(dir, event.Name))
					if len(messages) == 0 {
						lock.Unlock()
						continue
					}
					localeMessages[pathx.FileNoExt(event.Name)] = messages
					lock.Unlock()
					log.Infow("locale file changed", zap.String("file", event.Name), zap.Any("messages", messages))
				}
			case err, ok := <-watch.Errors:
				if !ok {
					return
				}
				log.Errorw("locale file watch error", zap.Error(err))
			}
		}
	}()

}

func loadLocaleFile(file string) map[string]string {
	path := filepath.Clean(file)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	messages := make(map[string]string)
	if err = json.Unmarshal(data, &messages); err != nil {
		return nil
	}
	return messages
}
