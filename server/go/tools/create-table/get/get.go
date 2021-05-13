package get

import (
	"github.com/liov/hoper/v2/tiga/initialize"
	"github.com/liov/hoper/v2/tiga/initialize/inject_dao"
	"gorm.io/gorm"
)

type Dao struct {
	GORMDB *gorm.DB `config:"database"`
}

type Config struct {
	GORMDB inject_dao.DatabaseConfig
}

func (*Config) Init() {}
func (*Dao) Init()    {}
func (*Dao) Close()   {}

var config = Config{}
var dao = Dao{}

func GetDB() *gorm.DB {
	if dao.GORMDB != nil {
		return dao.GORMDB
	}
	initialize.Start(&config, &dao)
	return dao.GORMDB
}
