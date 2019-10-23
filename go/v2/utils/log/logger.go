package log

import (
	"fmt"
	"log"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.SugaredLogger
}

func (l *Logger) Print(args ...interface{}) {
	l.Info(args...)
}

func GetLogger() *Logger {
	return logger
}

type LoggerInfo struct {
	Development bool
	OutLevel    zapcore.Level
	OutputPaths []string //日志文件路径
	ErrOutputPaths []string
	ServiceName string   //系统名称namespace.service
	LoggerCall
}

//初始化日志对象
func (lf *LoggerInfo) newLogger() *Logger {
	return &Logger{
		lf.initLogger(),
	}
}

func (lf *LoggerInfo) NewLogger() {
	logger.SugaredLogger = lf.initLogger()
	Default = lf.noCaller()
}


//构建日志对象基本信息
func (lf *LoggerInfo) initConfig() (*zap.Config,zap.Option) {
	config := zap.NewProductionConfig()
	config.Level = zap.NewAtomicLevelAt(lf.OutLevel)
	config.EncoderConfig = zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    lf.ServiceName,
		CallerKey:  "caller",
		MessageKey: "msg",
		//StacktraceKey: "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	var hook zap.Option
	//hook日志初始化

	if len(lf.HookURL) != 0 {
		lf.AtMan = "@" + strings.ReplaceAll(lf.AtMan, ",", "@")

		hook = zap.Hooks(func(e zapcore.Entry) error {
			return lf.Fire(&e)
		})
	} else {
		hook = zap.Hooks(func(i zapcore.Entry) error {
			return nil
		})
	}

	if lf.Development {
		config.Encoding = "console"
		config.Sampling = nil
		config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		hook = zap.Development()
	} else {
		config.EncoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		//系统名称
		config.InitialFields = map[string]interface{}{"source": lf.ServiceName}
		//输出文件
		config.OutputPaths = lf.OutputPaths
	}
	return &config,hook
}


//构建日志对象基本信息
func (lf *LoggerInfo) initLogger() *zap.SugaredLogger {
	config,hook:=lf.initConfig()
	logger, err := config.Build(zap.AddCallerSkip(1), hook)
	if err != nil {
	log.Fatalf("open file error :%s", err.Error())
	}
	return logger.Sugar()
}

func (l *Logger) ReportLog(interfaceName string, msg ...interface{}) {
	l.With("fields", map[string]interface{}{
		"interface": interfaceName,
	}).Warn(msg...)
}

var logger *Logger = (&LoggerInfo{Development: true}).newLogger()

func Sync() {
	logger.Sync()
}

func Print(v ...interface{})  {

}

func Debug(v ...interface{}) {
	logger.Debug(v...)
}

func Info(v ...interface{}) {
	logger.Info(v...)
}

func Warn(format string, v ...interface{}) {
	logger.Warn(v...)
}

func Error(v ...interface{}) {
	logger.Error(v...)
}

func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Fatal(msg)
}

func Debugf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Debug(msg)
}

func Infof(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Info(msg)
}

func Warnf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Warn(msg)
}

func Errorf(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	logger.Error(msg)
}