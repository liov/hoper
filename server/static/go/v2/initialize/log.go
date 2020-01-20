package initialize

import (
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
)

type LogConfig log.Config

func (conf *LogConfig) Generate() *log.Config {
	return (*log.Config)(conf)
}

func (init *Init) P1Log() {
	conf := &LogConfig{}
	if exist := reflect3.GetFieldValue(init.conf, conf); !exist {
		return
	}
	logConf := conf.Generate()
	logConf.Development = init.Env == DEVELOPMENT
	logConf.ModuleName = init.Module
	logConf.SetLogger()
}
