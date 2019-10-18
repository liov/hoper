package initialize

import (
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/liov/hoper/go/v2/utils/log"
)

func (i *Init) P1Log(conf interface{}) {
	logConf:=LogConfig{}
	if exist := reflect3.GetExpectTypeValue(conf,&logConf);!exist{
		return
	}
	(&log.LoggerInfo{
		Development: i.Env != PRODUCT,
		OutLevel:    logConf.Level,
		LogFilePath: logConf.FilePath,
	}).NewLogger()
}