package initialize

import (
	"github.com/liov/hoper/go/v2/utils/h_reflect"
	"github.com/liov/hoper/go/v2/utils/log"
)

func (i *Init) P1Log(conf interface{}) {
	logConf:=LogConfig{}
	if exist := h_reflect.GetExpectTypeValue(conf,&logConf);!exist{
		return
	}
	(&log.LoggerInfo{
		Development: i.Env != PRODUCT,
		OutLevel:    logConf.Level,
		LogFilePath: logConf.FilePath,
	}).NewLogger()
}