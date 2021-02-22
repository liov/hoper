package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SugaredLogger struct {
	zap.SugaredLogger
	Level zapcore.Level
}

// 兼容grpclog
func (l *SugaredLogger) Infoln(v ...interface{}) {
	l.Info(v...)
}

func (l *SugaredLogger) Warning(v ...interface{}) {
	l.Warn(v...)
}

func (l *SugaredLogger) Warningln(v ...interface{}) {
	l.Warn(v...)
}

func (l *SugaredLogger) Warningf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

func (l *SugaredLogger) Errorln(v ...interface{}) {
	l.Warn(v...)
}

func (l *SugaredLogger) Fatalln(v ...interface{}) {
	l.Warn(v...)
}

func (l *SugaredLogger) V(level int) bool {
	if level == 3 {
		level = 5
	}
	return int(l.Level) >= level
}
