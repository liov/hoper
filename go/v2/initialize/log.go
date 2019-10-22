package initialize

import (
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/liov/hoper/go/v2/utils/log"
)

func (i *Init) P1Log() {
	logConf:=LogConfig{}
	if exist := reflect3.GetFieldValue(i.conf,&logConf);!exist{
		return
	}
	(&log.LoggerInfo{
		Development: i.Env != PRODUCT,
		OutLevel:    logConf.Level,
		LogFilePath: logConf.FilePath,
	}).NewLogger()
}