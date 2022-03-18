package get_db

import (
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize"
	"github.com/actliboy/hoper/server/go/lib/tiga/initialize/db"
	"gorm.io/gorm"
)

type Dao struct {
	Hoper db.DB
}

type Config struct {
}

func (*Config) Init() {}
func (*Dao) Init()    {}
func (*Dao) Close()   {}

var config = Config{}
var dao = Dao{}

func GetDB() *gorm.DB {
	if dao.Hoper.DB != nil {
		return dao.Hoper.DB
	}
	initialize.Start(&config, &dao)
	return dao.Hoper.DB
}
