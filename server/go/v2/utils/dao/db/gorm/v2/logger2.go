package v2

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type SQLLogger struct {
	logger *zap.Logger
}

func New2(loger *zap.Logger, conf *logger.Config) logger.Interface {
	if conf == nil {
		config = &logger.Config{LogLevel: logger.Silent}
	}
	config = conf
	loger.Core().Enabled(zapcore.Level(4 - conf.LogLevel))
	return &SQLLogger{logger: loger}
}

// LogMode log mode
func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
	l.logger.Core().Enabled(zapcore.Level(4 - level))
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
		line := utils.FileWithLineNum()
		switch {
		case err != nil:
			if rows == -1 {
				l.logger.Error(line, zap.Error(err), elapsedms, sqlField)
			} else {
				l.logger.Error(line, zap.Error(err), elapsedms, rowsField, sqlField)
			}
		case elapsed > config.SlowThreshold && config.SlowThreshold != 0:

			slowLog := zap.String("slowLog", fmt.Sprintf("SLOW SQL >= %v", config.SlowThreshold))
			if rows == -1 {
				l.logger.Warn(line, slowLog, elapsedms, sqlField)
			} else {
				l.logger.Warn(line, slowLog, elapsedms, rowsField, sqlField)
			}
		default:
			if rows == -1 {
				l.logger.Info(line, elapsedms, sqlField)
			} else {
				l.logger.Info(line, elapsedms, rowsField, sqlField)
			}
		}
	}
}
