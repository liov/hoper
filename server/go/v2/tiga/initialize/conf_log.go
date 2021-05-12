package initialize

import (
	"github.com/liov/hoper/go/v2/utils/log"
)

type LogConfig log.Config

func (conf *LogConfig) Init() {
	logConf := (*log.Config)(conf)
	logConf.Development = InitConfig.Env != PRODUCT
	logConf.ModuleName = InitConfig.Module
	logConf.SetLogger()
}
