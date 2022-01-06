package get_db

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/inject_dao"
	"gorm.io/gorm"
)

type Dao struct {
	Hoper *gorm.DB `config:"database"`
}

type Config struct {
	Hoper inject_dao.DatabaseConfig
}

func (*Config) Init() {}
func (*Dao) Init()    {}
func (*Dao) Close()   {}

var config = Config{}
var dao = Dao{}

func GetDB() *gorm.DB {
	if dao.Hoper != nil {
		return dao.Hoper
	}
	initialize.Start(&config, &dao)
	return dao.Hoper
}
