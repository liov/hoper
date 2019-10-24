package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/liov/hoper/go/v2/utils/log/output"
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
	Product     bool
	Kafka       *output.LogKafka
	OutLevel    zapcore.Level
	OutputPaths []string //日志文件路径
	ModuleName  string   //系统名称namespace.service
	LoggerCall
}

//初始化日志对象
func (lf *LoggerInfo) NewLogger() *Logger {
	return &Logger{
		lf.initConfig(true).Sugar(),
	}
}

func (lf *LoggerInfo) SetLogger() {
	logger.SugaredLogger = lf.initConfig(true).Sugar()
}

//构建日志对象基本信息
func (lf *LoggerInfo) initConfig(skip bool) *zap.Logger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    lf.ModuleName,
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
	var hooks []zap.Option
	//hook日志初始化
	if len(lf.HookURL) != 0 {
		lf.AtMan = "@" + strings.ReplaceAll(lf.AtMan, ",", "@")

		hooks = append(hooks, zap.Hooks(func(e zapcore.Entry) error {
			return lf.Fire(&e)
		}))
	} else {
		hooks = append(hooks, zap.Hooks(func(i zapcore.Entry) error {
			return nil
		}))
	}

	encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	sink, _, err := zap.Open(lf.OutputPaths...)
	if err != nil {
		log.Fatal(err)
	}

	core:= zapcore.NewCore(jsonEncoder,sink,lf.OutLevel)
	//系统名称
	hooks= append(hooks,zap.Fields(zap.Any("module", lf.ModuleName)))

	if !lf.Product {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		core = zapcore.NewTee(core,zapcore.NewCore(consoleEncoder,zapcore.AddSync(os.Stderr),lf.OutLevel))
		hooks = append(hooks,zap.Development())
	}
	if skip {
		hooks = append(hooks,zap.AddCaller(),zap.AddCallerSkip(1))
	}

	logger :=  zap.New(core, hooks..., )

	return logger
}

func getWriter(filename string) io.Writer {
	// 生成rotatelogs的Logger 实际生成的文件名 demo.log.YYmmddHH
	// demo.log是指向最新日志的链接
	// 保存7天内的日志，每1小时(整点)分割一次日志
	hook, err := rotatelogs.New(
		filename+"_%Y%m%d%H", // 没有使用go风格反人类的format格式
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(time.Hour*24*7),
		rotatelogs.WithRotationTime(time.Hour),
	)

	if err != nil {
		panic(err)
	}
	return hook
}


func (l *Logger) ReportLog(interfaceName string, msg ...interface{}) {
	l.With("fields", map[string]interface{}{
		"interface": interfaceName,
	}).Warn(msg...)
}

var logger *Logger = (&LoggerInfo{Product: true}).NewLogger()

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