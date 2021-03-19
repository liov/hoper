package gormi

import (
	v2 "github.com/liov/hoper/go/v2/utils/dao/db/gorm/v2"
	"github.com/liov/hoper/go/v2/utils/log"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func GetDB(db *gorm.DB, log *log.Logger,conf *logger.Config) *gorm.DB {
	return db.Session(&gorm.Session{
		Logger: &v2.SQLLogger{Logger: log.Logger,
			Config: conf,
		}})
}

