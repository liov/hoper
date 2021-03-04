package gormi

import "gorm.io/gorm/logger"

type GORMConfig struct {
	//*gorm.Config
	Logger *logger.Config
}
