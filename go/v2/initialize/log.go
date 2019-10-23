package initialize

import (
	"github.com/liov/hoper/go/v2/utils/reflect3"
	"github.com/liov/hoper/go/v2/utils/log"
	"go.uber.org/zap/zapcore"
)

func (i *Init) P1Log() {
	logConf:=LogConfig{}
	if exist := reflect3.GetFieldValue(i.conf,&logConf);!exist{
		return
	}
	(&log.LoggerInfo{
		Development: i.Env != PRODUCT,
		OutLevel:    zapcore.Level(logConf.Level),
		OutputPaths: logConf.OutputPaths,
		ErrOutputPaths:logConf.ErrOutputPaths,
	}).NewLogger()
}