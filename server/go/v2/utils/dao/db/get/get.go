package get

import (
	"github.com/liov/hoper/go/v2/initialize"
	"gorm.io/gorm"
)

type Dao struct {
	GORMDB *gorm.DB
}

type Config struct {
	Database initialize.DatabaseConfig
}

func (*Config) Custom() {}
func (*Dao) Custom()    {}
func (*Dao) Close()     {}

var config = Config{}
var dao = Dao{}

func GetDB() *gorm.DB {
	if dao.GORMDB != nil {
		return dao.GORMDB
	}
	initialize.Start(&config, &dao)
	return dao.GORMDB
}
