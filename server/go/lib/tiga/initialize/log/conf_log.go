package log

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

type LogConfig log.Config

func (conf *LogConfig) Init() {
	logConf := (*log.Config)(conf)
	logConf.Development = initialize.InitConfig.Env != initialize.PRODUCT
	logConf.ModuleName = initialize.InitConfig.Module
	log.SetDefaultLogger(logConf)
}
