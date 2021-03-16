package gormi

import (
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type GORMConfig struct {
	Config gorm.Config
	Logger logger.Config
}

func (c *GORMConfig) Custom()    {
	c.Logger.SlowThreshold *= time.Millisecond
}