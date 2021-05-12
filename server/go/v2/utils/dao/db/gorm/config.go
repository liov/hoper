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

func (c *GORMConfig) Init() {
	if c.Logger.SlowThreshold < 10*time.Millisecond {
		c.Logger.SlowThreshold *= time.Millisecond
	}
}
