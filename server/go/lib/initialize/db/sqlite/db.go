package sqlite

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	pkdb "github.com/actliboy/hoper/server/go/lib/initialize/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"runtime"
)

type DatabaseConfig pkdb.DatabaseConfig

func (conf *DatabaseConfig) Generate() any {
	return conf.generate()
}

func (conf *DatabaseConfig) Init() {
	(*pkdb.DatabaseConfig)(conf).Init()
}

func (conf *DatabaseConfig) generate() *gorm.DB {
	url := "/data/db/sqlite/" + conf.Database + ".db"
	if runtime.GOOS == "windows" {
		url = ".." + url
	}
	return (*pkdb.DatabaseConfig)(conf).Generate(sqlite.Open(url))
}

type DB pkdb.DB

func (db *DB) Config() initialize.Generate {
	return (*DatabaseConfig)(&db.Conf)
}

func (db *DB) SetEntity(entity interface{}) {
	if gormdb, ok := entity.(*gorm.DB); ok {
		db.DB = gormdb
	}
}

func (db *DB) Close() error {
	return nil
}
