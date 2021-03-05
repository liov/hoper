package gormi

import (
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GORMConfig struct {
	Config gorm.Config
	Logger logger.Config
}
