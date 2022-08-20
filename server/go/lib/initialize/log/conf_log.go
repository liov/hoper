package log

import (
	initialize2 "github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/utils/log"
)

type LogConfig log.Config

func (conf *LogConfig) Init() {
	logConf := (*log.Config)(conf)
	logConf.Development = initialize2.InitConfig.Env != initialize2.PRODUCT
	logConf.ModuleName = initialize2.InitConfig.Module
	log.SetDefaultLogger(logConf)
}

func (conf *LogConfig) generate() *log.Logger {

	return (*log.Config)(conf).NewLogger()
}

func (conf *LogConfig) Generate() interface{} {
	return conf.generate()
}

type Logger struct {
	*log.Logger `init:"entity"`
	Conf        LogConfig `init:"config"`
}

func (l *Logger) Config() initialize2.Generate {
	return &l.Conf
}

func (l *Logger) SetEntity(entity interface{}) {
	if logger, ok := entity.(*log.Logger); ok {
		l.Logger = logger
	}
}

func (l *Logger) Close() error {
	return l.Logger.Sync()
}
