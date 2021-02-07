package log

import (
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
	*zap.SugaredLogger
	Logger *zap.Logger
}

func GetLogger() *Logger {
	return Default
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
	return &Logger{logger.Sugar(), logger}
}

func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	logger :=*l
	logger.Logger = l.Logger.WithOptions(opts...)
	logger.SugaredLogger = logger.Logger.Sugar()
	return &logger
}

func (l *Logger) WithFields(fields ...zap.Field) *Logger {
	logger :=*l
	logger.Logger = l.Logger.With(fields...)
	logger.SugaredLogger = logger.Logger.Sugar()
	return &logger
}

var Default = (&Config{Development: true, Caller: true, Level: -1}).NewLogger()
var CallOne = Default.WithOptions(zap.AddCallerSkip(1))

func init() {
	output.RegisterSink()
}

func (lf *Config) SetLogger() {
	Default = lf.NewLogger()
	CallOne = Default.WithOptions(zap.AddCallerSkip(1))
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
		encoderConfig.NameKey =  "module"
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
	if lf.ModuleName != ""{
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

func Print(v ...interface{}) {
	CallOne.Print(v...)
}

func Debug(v ...interface{}) {
	CallOne.Debug(v...)
}

func Info(v ...interface{}) {
	CallOne.Info(v...)
}

func Warn(format string, v ...interface{}) {
	CallOne.Warn(v...)
}

func Error(v ...interface{}) {
	CallOne.Error(v...)
}

func Panic(v ...interface{}) {
	CallOne.Panic(v...)
}

func Fatal(v ...interface{}) {
	CallOne.Fatal(v...)
}

func Printf(format string, v ...interface{}) {
	CallOne.Printf(format, v...)
}

func Debugf(format string, v ...interface{}) {
	CallOne.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	CallOne.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	CallOne.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	CallOne.Errorf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	CallOne.Panicf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	CallOne.Fatalf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	CallOne.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	CallOne.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	CallOne.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	CallOne.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	CallOne.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	CallOne.Fatalw(msg, keysAndValues...)
}

// 兼容gormv1
func (l *Logger) Printf(format string, args ...interface{}) {
	l.Infof(format, args...)
}

func (l *Logger) Print(args ...interface{}) {
	l.Info(args...)
}

// 兼容grpclog
func (l *Logger) Infoln(v ...interface{}) {
	l.Info(v...)
}

func (l *Logger) Warning(v ...interface{}) {
	l.Warn(v...)
}

func (l *Logger) Warningln(v ...interface{}) {
	l.Warn(v...)
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

func (l *Logger) Errorln(v ...interface{}) {
	l.Warn(v...)
}

func (l *Logger) Fatalln(v ...interface{}) {
	l.Warn(v...)
}

func (l *Logger) V(level int) bool {
	if level == 3 {
		level = 5
	}
	return l.Logger.Core().Enabled(zapcore.Level(level))
}
