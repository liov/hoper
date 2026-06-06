package global

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hopeio/initialize"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.opentelemetry.io/otel"
	"golang.org/x/text/language"
)

const ScopeName = "github.com/liov/hoper/server/go/user"

var (
	Global         = initialize.NewGlobal[*config, *dao]()
	Dao    *dao    = Global.Dao
	Conf   *config = Global.Config

	Tracer = otel.Tracer(ScopeName)
	Meter  = otel.Meter(ScopeName)

	LocalizerMap = make(map[string]*i18n.Localizer)
)

func init() {
	bundle := i18n.NewBundle(language.Chinese)
	// 1. 加载翻译文件
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	dir := Conf.Locale.Dir
	files, err := os.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		_, err = bundle.LoadMessageFile(filepath.Join(dir, file.Name()))
		if err != nil {
			log.Fatal(err)
		}
		locale := strings.Split(file.Name(), ".")[0]
		LocalizerMap[locale] = i18n.NewLocalizer(bundle, locale)
	}
}
