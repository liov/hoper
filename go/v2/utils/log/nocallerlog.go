package log

import (
	"log"

	"go.uber.org/zap"
)

var  Default = (&LoggerInfo{Development: true}).noCaller()

func (lf *LoggerInfo) noCaller() *zap.SugaredLogger {
	config,hook:=lf.initConfig()
	config.DisableCaller = true
	logger, err := config.Build(hook)
	if err != nil {
		log.Fatalf("open file error :%s", err.Error())
	}
	return logger.Sugar()
}