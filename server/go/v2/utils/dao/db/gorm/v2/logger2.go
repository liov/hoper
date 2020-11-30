package v2

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

type SQLLogger struct {
	logger *zap.Logger
	logger.LogLevel
}

func Config(conf *logger.Config) {
	if conf == nil {
		config = &logger.Config{}
	}
	config = conf
}

// LogMode log mode
func (l *SQLLogger) LogMode(level logger.LogLevel) logger.Interface {
	newlogger := *l
	newlogger.LogLevel = level
	return &newlogger
}

// Info print info
func (l *SQLLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger.Info(msg)
	}
}

// Warn print warn messages
func (l *SQLLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger.Warn(msg)
	}
}

// Error print error messages
func (l *SQLLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger.Error(msg)
	}
}

// Trace print sql message
func (l *SQLLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.LogLevel > 0 {
		elapsed := time.Since(begin)
		elapsedms := zap.Float64("elapsedms", float64(elapsed.Nanoseconds())/1e6)
		sql, rows := fc()
		sqlField := zap.String("sql", sql)
		rowsField := zap.Int64("rows", rows)
		line := utils.FileWithLineNum()
		switch {
		case err != nil && l.LogLevel >= logger.Error:
			if rows == -1 {
				l.logger.Error(line, zap.Error(err), elapsedms, sqlField)
			} else {
				l.logger.Error(line, zap.Error(err), elapsedms, rowsField, sqlField)
			}
		case elapsed > config.SlowThreshold && config.SlowThreshold != 0 && l.LogLevel >= logger.Warn:

			slowLog := zap.String("slowLog", fmt.Sprintf("SLOW SQL >= %v", config.SlowThreshold))
			if rows == -1 {
				l.logger.Warn(line, slowLog, elapsedms, sqlField)
			} else {
				l.logger.Warn(line, slowLog, elapsedms, rowsField, sqlField)
			}
		case l.LogLevel >= logger.Info:

			if rows == -1 {
				l.logger.Info(line, elapsedms, sqlField)
			} else {
				l.logger.Info(line, elapsedms, rowsField, sqlField)
			}
		}
	}
}
