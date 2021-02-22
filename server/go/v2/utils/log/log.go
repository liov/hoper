package log

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"

	"github.com/liov/hoper/go/v2/utils/log/output"
	"github.com/liov/hoper/go/v2/utils/net"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
	Level zapcore.Level
}

type Config struct {
	Development bool
	Caller      bool
	Level       zapcore.Level
	OutputPaths map[string][]string
	ModuleName  string //系统名称namespace.service
}

//初始化日志对象
func (lf *Config) NewLogger() *Logger {
	logger := lf.initLogger().
		With(
			zap.String("source", neti.GetIP()),
		)
	return &Logger{logger, lf.Level}
}

func (l *Logger) Sugar() *SugaredLogger {
	return &SugaredLogger{*l.Logger.Sugar(), l.Level}
}

func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	return &Logger{l.Logger.WithOptions(opts...), l.Level}
}

func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...), l.Level}
}

var (
	CallOne      *Logger
	Default      *Logger
	SugarCallOne *SugaredLogger
	Sugar        *SugaredLogger
)

func (lf *Config) SetLogger() {
	Default = lf.NewLogger()
	CallOne = Default.WithOptions(zap.AddCallerSkip(1))
	Sugar = Default.Sugar()
	SugarCallOne = Default.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func init() {
	output.RegisterSink()
	(&Config{Development: true, Caller: true, Level: -1}).SetLogger()
}

//构建日志对象基本信息
func (lf *Config) initLogger() *zap.Logger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		CallerKey:     "caller",
		FunctionKey:   "func",
		MessageKey:    "msg",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(runtime.FuncForPC(caller.PC).Name() + ` ` + caller.TrimmedPath())
		},
	}
	if lf.ModuleName != "" {
		encoderConfig.NameKey = "module"
	}
	if lf.Development {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
	}

	var consoleEncoder, jsonEncoder zapcore.Encoder
	var cores []zapcore.Core

	if len(lf.OutputPaths["console"]) > 0 {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		sink, _, err := zap.Open(lf.OutputPaths["console"]...)
		if err != nil {
			log.Fatal(err)
		}
		cores = append(cores, zapcore.NewCore(consoleEncoder, sink, lf.Level))
	}

	if len(lf.OutputPaths["json"]) > 0 {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
		jsonEncoder = zapcore.NewJSONEncoder(encoderConfig)
		sink, _, err := zap.Open(lf.OutputPaths["json"]...)
		if err != nil {
			log.Fatal(err)
		}
		cores = append(cores, zapcore.NewCore(jsonEncoder, sink, lf.Level))
	}
	//如果没有设置输出，默认控制台
	if len(cores) == 0 {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), lf.Level))
	}
	core := zapcore.NewTee(cores...)

	logger := zap.New(core, lf.hook()...)
	if lf.ModuleName != "" {
		logger = logger.Named(lf.ModuleName)
	}
	return logger
}

func (lf *Config) hook() []zap.Option {
	var hooks []zap.Option

	if lf.Development {
		hooks = append(hooks, zap.Development(), zap.AddStacktrace(zapcore.DPanicLevel))
	}
	if lf.Caller {
		hooks = append(hooks, zap.AddCaller())
	}
	return hooks
}

func Sync() {
	Default.Sync()
	CallOne.Sync()
}

func Print(args ...interface{}) {
	CallOne.Print(args...)
}

func Debug(args ...interface{}) {
	CallOne.Debug(fmt.Sprint(args...))
}

func Info(args ...interface{}) {
	CallOne.Info(fmt.Sprint(args...))
}

func Warn(args ...interface{}) {
	CallOne.Warn(fmt.Sprint(args...))
}

func Error(args ...interface{}) {
	CallOne.Error(fmt.Sprint(args...))
}

func Panic(args ...interface{}) {
	CallOne.Panic(fmt.Sprint(args...))
}

func Fatal(args ...interface{}) {
	CallOne.Fatal(fmt.Sprint(args...))
}

func Printf(template string, args ...interface{}) {
	CallOne.Print(fmt.Sprintf(template, args...))
}

func Debugf(template string, args ...interface{}) {
	CallOne.Debug(fmt.Sprintf(template, args...))
}

func Infof(template string, args ...interface{}) {
	CallOne.Info(fmt.Sprintf(template, args...))
}

func Warnf(template string, args ...interface{}) {
	CallOne.Warn(fmt.Sprintf(template, args...))
}

func Errorf(template string, args ...interface{}) {
	CallOne.Error(fmt.Sprintf(template, args...))
}

func Panicf(template string, args ...interface{}) {
	CallOne.Panic(fmt.Sprintf(template, args...))
}

func Fatalf(template string, args ...interface{}) {
	CallOne.Fatal(fmt.Sprintf(template, args...))
}

// 兼容gormv1
func (l *Logger) Printf(template string, args ...interface{}) {
	l.Info(fmt.Sprintf(template, args...))
}

func (l *Logger) Print(args ...interface{}) {
	l.Info(fmt.Sprint(args...))
}
