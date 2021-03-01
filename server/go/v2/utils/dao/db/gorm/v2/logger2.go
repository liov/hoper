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
)

type SQLLogger struct {
	logger *zap.Logger
	*logger.Config
}

func New2(loger *zap.Logger, conf *logger.Config) logger.Interface {
	if conf == nil {
		conf = &logger.Config{LogLevel: logger.Silent}
	}
	loger.Core().Enabled(zapcore.Level(4 - conf.LogLevel))
	return &SQLLogger{logger: loger, Config: conf}
}

// LogMode log mode
func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logger.Core().Enabled(zapcore.Level(4 - level))
	l.LogLevel = level
	return l
}

// Info print info
func (l *SQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Info(fmt.Sprintf(msg, data...))
}

// Warn print warn messages
func (l *SQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Warn(fmt.Sprintf(msg, data...))
}

// Error print error messages
func (l *SQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logger.Error(fmt.Sprintf(msg, data...))
}

// Trace print sql message
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if config.LogLevel > 0 {
		elapsed := time.Since(begin)
		elapsedms := zap.Float64("elapsedms", float64(elapsed.Nanoseconds())/1e6)
		sql, rows := fc()
		sqlField := zap.String("sql", sql)
		rowsField := zap.Int64("rows", rows)
		fields := make([]zap.Field, 3, 4)
		fields[0], fields[1], fields[2] = elapsedms, sqlField, rowsField
		line := utils.FileWithLineNum()
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			fields = append(fields, zap.Error(err))
			l.logger.Error(line, fields...)
		case elapsed > config.SlowThreshold && config.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
			fields = append(fields, zap.String("slowLog", fmt.Sprintf("SLOW SQL >= %v", config.SlowThreshold)))
			l.logger.Warn(line, fields...)
		case l.LogLevel >= logger.Info:
			l.logger.Info(line, fields...)
		}
	}
}
