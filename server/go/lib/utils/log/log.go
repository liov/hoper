package log

import (
	"fmt"
	"github.com/liov/hoper/server/go/lib/utils/log/output"
	neti "github.com/liov/hoper/server/go/lib/utils/net"
	osi "github.com/liov/hoper/server/go/lib/utils/os"
	"log"
	"os"
	"strconv"
	"time"

	"go.uber.org/multierr"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	output.RegisterSink()
	SetDefaultLogger(&Config{Development: true, Caller: true, Level: zapcore.DebugLevel})
}

var (
	Default     *Logger
	skipLoggers = make([]*Logger, 10)
)

func SetDefaultLogger(lf *Config) {
	Default = lf.NewLogger()
}

func GetSkipLogger(skip int) *Logger {
	if skip > 10 {
		panic("最大不超过10")
	}
	if skipLoggers[skip] == nil {
		skipLoggers[skip] = Default.AddSkip(skip)
	}
	return skipLoggers[skip]
}

type Logger struct {
	*zap.Logger
}

type ZapConfig zap.Config

type Config struct {
	Development      bool
	Caller           bool
	Level            zapcore.Level
	OutputPaths      map[string][]string
	ErrorOutputPaths map[string][]string
	ModuleName       string //系统名称namespace.service
}

// 初始化日志对象
func (lf *Config) NewLogger(cores ...zapcore.Core) *Logger {
	logger := lf.initLogger(cores...).
		With(
			zap.String("hostname", osi.Hostname()),
			zap.String("ip", neti.GetIP()),
		)
	return &Logger{logger}
}

// Named adds a sub-scope to the logger's name. See Logger.Named for details.
func (l *Logger) Named(name string) *Logger {
	return &Logger{l.Logger.Named(name)}
}

func (l *Logger) WithOptions(opts ...zap.Option) *Logger {
	return &Logger{l.Logger.WithOptions(opts...)}
}

func (l *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{l.Logger.With(fields...)}
}

func (l *Logger) Sugar() *zap.SugaredLogger {
	l.WithOptions(zap.AddCallerSkip(-1))
	return l.Logger.Sugar()
}

func (l *Logger) AddCore(newCore zapcore.Core) *Logger {
	return l.WithOptions(zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.NewTee(core, newCore)
	}))
}

func (l *Logger) AddSkip(skip int) *Logger {
	return &Logger{l.Logger.WithOptions(zap.AddCallerSkip(skip))}
}

// 构建日志对象基本信息
func (lf *Config) initLogger(cores ...zapcore.Core) *zap.Logger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:       "time",
		LevelKey:      "level",
		CallerKey:     "caller",
		FunctionKey:   "func",
		MessageKey:    "msg",
		StacktraceKey: "stack",
		LineEnding:    zapcore.DefaultLineEnding,
		EncodeLevel:   zapcore.CapitalColorLevelEncoder,
		EncodeTime: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006/01/02 15:04:05.000"))
		},
		EncodeDuration: func(d time.Duration, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(strconv.FormatInt(d.Nanoseconds()/1e6, 10) + "ms")
		},
		EncodeCaller: func(caller zapcore.EntryCaller, encoder zapcore.PrimitiveArrayEncoder) {
			encoder.AppendString(caller.TrimmedPath())
		},
		ConsoleSeparator: "\t",
	}
	if lf.ModuleName != "" {
		encoderConfig.NameKey = "module"
	}
	if lf.Development {
		encoderConfig.EncodeCaller = zapcore.FullCallerEncoder
		//encoderConfig.ConsoleSeparator = "\n"
	}

	var consoleEncoder, jsonEncoder zapcore.Encoder

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
	hooks = append(hooks, zap.AddCallerSkip(1))
	if lf.Caller {
		hooks = append(hooks, zap.AddCaller())
	}
	return hooks
}

func Sync() {
	Default.Sync()
}

func Print(args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Debug(args ...interface{}) {
	if ce := Default.Check(zap.DebugLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Info(args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Warn(args ...interface{}) {
	if ce := Default.Check(zap.WarnLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Error(args ...interface{}) {
	if ce := Default.Check(zap.ErrorLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Panic(args ...interface{}) {
	if ce := Default.Check(zap.PanicLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Fatal(args ...interface{}) {
	if ce := Default.Check(zap.FatalLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func Printf(template string, args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Debugf(template string, args ...interface{}) {
	if ce := Default.Check(zap.DebugLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Infof(template string, args ...interface{}) {
	if ce := Default.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Warnf(template string, args ...interface{}) {
	if ce := Default.Check(zap.WarnLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Errorf(template string, args ...interface{}) {
	if ce := Default.Check(zap.ErrorLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Panicf(template string, args ...interface{}) {
	if ce := Default.Check(zap.PanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func Fatalf(template string, args ...interface{}) {
	if ce := Default.Check(zap.FatalLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Printf(template string, args ...interface{}) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// 兼容gormv1
func (l *Logger) Print(args ...interface{}) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(args ...interface{}) {
	if ce := l.Check(zap.DebugLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(args ...interface{}) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(args ...interface{}) {
	if ce := l.Check(zap.WarnLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(args ...interface{}) {
	if ce := l.Check(zap.ErrorLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// DPanic uses fmt.Sprint to construct and log a message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Logger) DPanic(args ...interface{}) {
	if ce := l.Check(zap.DPanicLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(args ...interface{}) {
	if ce := l.Check(zap.PanicLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Logger) Fatal(args ...interface{}) {
	if ce := l.Check(zap.FatalLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Debugw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.DebugLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Infow(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.InfoLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Warnw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.WarnLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func (l *Logger) Errorw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.ErrorLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// DPanic logs a message at DPanicLevel. The message includes any fields
// passed at the log site, as well as any fields accumulated on the logger.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l *Logger) DPanicw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.DPanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l *Logger) Panicw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.PanicLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l *Logger) Fatalw(msg string, fields ...zap.Field) {
	if ce := l.Check(zap.FatalLevel, msg); ce != nil {
		ce.Write(fields...)
	}
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) Debugf(template string, args ...interface{}) {
	if ce := l.Check(zap.DebugLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(template string, args ...interface{}) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(template string, args ...interface{}) {
	if ce := l.Check(zap.WarnLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(template string, args ...interface{}) {
	if ce := l.Check(zap.ErrorLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the
// logger then panics. (See DPanicLevel for details.)
func (l *Logger) DPanicf(template string, args ...interface{}) {
	if ce := l.Check(zap.DPanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(template string, args ...interface{}) {
	if ce := l.Check(zap.PanicLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(template string, args ...interface{}) {
	if ce := l.Check(zap.FatalLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

// 兼容grpclog
func (l *Logger) Infoln(args ...interface{}) {
	if ce := l.Check(zap.InfoLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Warning(args ...interface{}) {
	if ce := l.Check(zap.WarnLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Warningln(args ...interface{}) {
	if ce := l.Check(zap.WarnLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Warningf(template string, args ...interface{}) {
	if ce := l.Check(zap.WarnLevel, fmt.Sprintf(template, args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Errorln(args ...interface{}) {
	if ce := l.Check(zap.ErrorLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

func (l *Logger) Fatalln(args ...interface{}) {
	if ce := l.Check(zap.FatalLevel, fmt.Sprint(args...)); ce != nil {
		ce.Write()
	}
}

// grpclog
func (l *Logger) V(level int) bool {
	level -= 2
	return l.Logger.Core().Enabled(zapcore.Level(level))
}

// sugar
const (
	_oddNumberErrMsg    = "Ignored key without a value."
	_nonStringKeyErrMsg = "Ignored key-value pairs with non-string keys."
)

func (l *Logger) log(lvl zapcore.Level, template string, fmtArgs []interface{}, context []interface{}) {
	// If logging at this level is completely disabled, skip the overhead of
	// string formatting.
	if lvl < zap.DPanicLevel && !l.Core().Enabled(lvl) {
		return
	}

	// Format with Sprint, Sprintf, or neither.
	msg := template
	if msg == "" && len(fmtArgs) > 0 {
		msg = fmt.Sprint(fmtArgs...)
	} else if msg != "" && len(fmtArgs) > 0 {
		msg = fmt.Sprintf(template, fmtArgs...)
	}

	if ce := l.Check(lvl, msg); ce != nil {
		ce.Write(l.sweetenFields(context)...)
	}
}

func (l *Logger) sweetenFields(args []interface{}) []zap.Field {
	if len(args) == 0 {
		return nil
	}

	// Allocate enough space for the worst case; if users pass only structured
	// fields, we shouldn't penalize them with extra allocations.
	fields := make([]zap.Field, 0, len(args))
	var invalid invalidPairs

	for i := 0; i < len(args); {
		// This is a strongly-typed field. Consume it and move on.
		if f, ok := args[i].(zap.Field); ok {
			fields = append(fields, f)
			i++
			continue
		}

		// Make sure this element isn't a dangling key.
		if i == len(args)-1 {
			l.DPanic(_oddNumberErrMsg, zap.Any("ignored", args[i]))
			break
		}

		// Consume this value and the next, treating them as a key-value pair. If the
		// key isn't a string, add this pair to the slice of invalid pairs.
		key, val := args[i], args[i+1]
		if keyStr, ok := key.(string); !ok {
			// Subsequent errors are likely, so allocate once up front.
			if cap(invalid) == 0 {
				invalid = make(invalidPairs, 0, len(args)/2)
			}
			invalid = append(invalid, invalidPair{i, key, val})
		} else {
			fields = append(fields, zap.Any(keyStr, val))
		}
		i += 2
	}

	// If we encountered any invalid key-value pairs, log an error.
	if len(invalid) > 0 {
		l.DPanic(_nonStringKeyErrMsg, zap.Array("invalid", invalid))
	}
	return fields
}

type invalidPair struct {
	position   int
	key, value interface{}
}

func (p invalidPair) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt64("position", int64(p.position))
	zap.Any("key", p.key).AddTo(enc)
	zap.Any("value", p.value).AddTo(enc)
	return nil
}

type invalidPairs []invalidPair

func (ps invalidPairs) MarshalLogArray(enc zapcore.ArrayEncoder) error {
	var err error
	for i := range ps {
		err = multierr.Append(err, enc.AppendObject(ps[i]))
	}
	return err
}
