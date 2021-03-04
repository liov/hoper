package v2

import (
	"context"
	"fmt"
	"time"
	logi "github.com/liov/hoper/go/v2/utils/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

var (
	DefaultV2 = New2(logi.Default.Logger, &logger.Config{
		SlowThreshold: 100 * time.Millisecond,
		LogLevel:      logger.Warn,
		Colorful:      true,
	})
	field = zap.String("db","gorm")
)

type SQLLogger struct {
	*zap.Logger
	*logger.Config
}

func New2(loger *zap.Logger, conf *logger.Config) logger.Interface {
	if conf == nil {
		conf = &logger.Config{LogLevel: logger.Warn}
	}
	loger.Core().Enabled(zapcore.Level(4 - conf.LogLevel))
	return &SQLLogger{Logger: loger, Config: conf}
}

// LogMode log mode
func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.Logger.Core().Enabled(zapcore.Level(4 - level))
	l.LogLevel = level
	return l
}

// Info print info
func (l *SQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Info(fmt.Sprintf(msg, data...),field)
}

// Warn print warn messages
func (l *SQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Warn(fmt.Sprintf(msg, data...),field)
}

// Error print error messages
func (l *SQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.Error(fmt.Sprintf(msg, data...),field)
}

// Trace print sql message 只有这里的context不是background,看了代码,也没用
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > logger.Silent {
		elapsed := time.Since(begin)
		elapsedms := zap.Float64("elapsedms", float64(elapsed.Nanoseconds())/1e6)
		sql, rows := fc()
		sqlField := zap.String("sql", sql)
		rowsField := zap.Int64("rows", rows)
		caller:=zap.String("caller",utils.FileWithLineNum())
		fields := []zap.Field{elapsedms, sqlField, rowsField,caller}
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			l.Logger.Error(err.Error(), fields...)
		case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			l.Logger.Warn(fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold), fields...)
		case l.LogLevel >= logger.Info:
			l.Logger.Info("", fields...)
		}
	}
}
