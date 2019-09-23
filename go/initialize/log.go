package initialize

import (
	"github.com/liov/hoper/go/user/internal/config"
	"github.com/liov/hoper/go/utls/log"
)

func (i *Init) P1Log() {
	var logConf = &config.Config.Log
	(&log.LoggerInfo{
		Development: i.Env != PRODUCT,
		OutLevel:    logConf.Level,
		LogFilePath: logConf.FilePath,
	}).NewLogger()
}
