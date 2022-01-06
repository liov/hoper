package initialize

import (
	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

type LogConfig log.Config

func (conf *LogConfig) Init() {
	logConf := (*log.Config)(conf)
	logConf.Development = InitConfig.Env != PRODUCT
	logConf.ModuleName = InitConfig.Module
	log.SetDefaultLogger(logConf)
}
