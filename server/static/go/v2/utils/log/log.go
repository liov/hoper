package log

import (
	"log"
	"os"
	"time"

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
	return Default
}

type Config struct {
	Development bool
	Skip        bool
	Level       zapcore.Level
	OutputPaths map[string][]string
	ModuleName  string //系统名称namespace.service
}

//初始化日志对象
func (lf *Config) NewLogger() *Logger {
	return &Logger{
		lf.initLogger().Sugar(),
	}
}

var Default *Logger = (&Config{Development: true, Skip: true, Level: -1}).NewLogger()
var NoCall = (&Config{Development: true}).NewLogger()
var CallTwo = Default.Desugar().WithOptions(zap.AddCallerSkip(3)).Sugar()

func init() {
	output.RegisterSink()
}

func (lf *Config) SetLogger() {
	Default.SugaredLogger = lf.initLogger().Sugar()
}

func (lf *Config) SetNoCall() {
	lf.Skip = false
	NoCall.SugaredLogger = lf.initLogger().Sugar()
}

//构建日志对象基本信息
func (lf *Config) initLogger() *zap.Logger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    lf.ModuleName,
		CallerKey:  "caller",
		MessageKey: "msg",
		//StacktraceKey: "stacktrace",
		LineEnding:  zapcore.DefaultLineEnding,
		EncodeLevel: zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006/01/02 - 15:04:05.000"))
		},
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
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

	if len(cores) == 0 {
		consoleEncoder = zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stderr), lf.Level))
	}
	core := zapcore.NewTee(cores...)

	logger := zap.New(core, lf.hook()...)

	return logger
}

func (lf *Config) hook() []zap.Option {
	var hooks []zap.Option
	//系统名称
	if len(lf.OutputPaths["json"]) > 0 && lf.ModuleName != "" {
		hooks = append(hooks, zap.Fields(zap.Any("module", lf.ModuleName)))
	}

	if lf.Development {
		hooks = append(hooks, zap.Development())
	}
	if lf.Skip {
		hooks = append(hooks, zap.AddCaller(), zap.AddCallerSkip(1))
	}
	return hooks
}

func Sync() {
	Default.Sync()
}

func Print(v ...interface{}) {

}

func Debug(v ...interface{}) {
	Default.Debug(v...)
}

func Info(v ...interface{}) {
	Default.Info(v...)
}

func Warn(format string, v ...interface{}) {
	Default.Warn(v...)
}

func Error(v ...interface{}) {
	Default.Error(v...)
}

func Panic(v ...interface{}) {
	Default.Panic(v...)
}

func Fatal(v ...interface{}) {
	Default.Fatal(v...)
}

func Debugf(format string, v ...interface{}) {
	Default.Debugf(format, v...)
}

func Infof(format string, v ...interface{}) {
	Default.Infof(format, v...)
}

func Warnf(format string, v ...interface{}) {
	Default.Warnf(format, v...)
}

func Errorf(format string, v ...interface{}) {
	Default.Errorf(format, v...)
}

func Panicf(format string, v ...interface{}) {
	Default.Panicf(format, v...)
}

func Fatalf(format string, v ...interface{}) {
	Default.Fatalf(format, v...)
}

func Debugw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	Default.Debugw(msg, keysAndValues...)
}

func Infow(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	Default.Infow(msg, keysAndValues...)
}

func Warnw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	Default.Warnw(msg, keysAndValues...)
}

func Errorw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	Default.Errorw(msg, keysAndValues...)
}

func Panicw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	Default.Panicw(msg, keysAndValues...)
}

func Fatalw(msg string, keysAndValues ...interface{}) {
	if len(keysAndValues)|0 != 0 {
		keysAndValues = append(keysAndValues, "")
	}
	Default.Fatalw(msg, keysAndValues...)
}
