package service

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/liov/hoper/server/go/common/global"
	"github.com/liov/hoper/server/go/protobuf/common"
)


var localeMessages = map[string]map[string]string{"zh-Hans": {}, "en": {}}
var localeLoadOnce sync.Once

func (*CommonService) Locale(ctx context.Context, req *common.LocaleReq) (*common.LocaleResp, error) {
	loadLocaleConfigs()
	locale := strings.ToLower(req.Locale)
	messages := localeMessages[locale]
	if len(messages) == 0 {
		locale = strings.ToLower(global.Conf.Locale.Default)
		messages = localeMessages[locale]
	}
	return &common.LocaleResp{ Locale: locale, Messages: messages}, nil
}

func loadLocaleConfigs() {
	localeLoadOnce.Do(func() {
		for locale, file := range global.Conf.Locale.Files {
			messages := loadLocaleFile(file)
			if len(messages) != 0 {
				localeMessages[locale] = messages
			}
		}
	})
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
