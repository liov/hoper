package gormi

import (
	"github.com/liov/hoper/v2/utils/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB(db *gorm.DB, log *log.Logger, conf *logger.Config) *gorm.DB {
	return db.Session(&gorm.Session{
		Logger: &SQLLogger{Logger: log.Logger,
			Config: conf,
		}})
}
