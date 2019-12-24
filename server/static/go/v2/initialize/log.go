package initialize

import (
	"github.com/liov/hoper/go/v2/utils/log"
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"go.uber.org/zap/zapcore"
)

type LogConfig struct {
	Level       int8
	Skip        bool
	OutputPaths map[string][]string
}

func (conf *LogConfig) Generate() *log.Config {
	return &log.Config{
		Development: true,
		ModuleName:  "",
		Skip:        conf.Skip,
		Level:       zapcore.Level(conf.Level),
		OutputPaths: conf.OutputPaths,
	}
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
