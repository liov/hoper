package initialize

import (
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/liov/hoper/go/v2/utils/log"
	"go.uber.org/zap/zapcore"
)

func (i *Init) P1Log() {
	conf := LogConfig{}
	if exist := reflect3.GetFieldValue(i.conf, &conf); !exist {
		return
	}
	(&log.Config{
		Development: i.Env == DEVELOPMENT,
		ModuleName:  i.Module,
		Skip:        conf.Skip,
		Level:       zapcore.Level(conf.Level),
		OutputPaths: conf.OutputPaths,
	}).SetLogger()
}
