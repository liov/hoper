package initialize

import (
	"reflect"

	"github.com/liov/hoper/go/v2/utils/log"
)

func (i *Init) P1Log(conf interface{}) {
	var logConf = (reflect.ValueOf(conf).Elem().FieldByName("Log").Interface()).(LogConfig)
	(&log.LoggerInfo{
		Development: i.Env != PRODUCT,
		OutLevel:    logConf.Level,
		LogFilePath: logConf.FilePath,
	}).NewLogger()
}
