package sqlite

import (
	pkdb "github.com/liov/hoper/server/go/lib/initialize/gormdb"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"runtime"
)

type DatabaseConfig pkdb.DatabaseConfig

func (conf *DatabaseConfig) Init() {
	(*pkdb.DatabaseConfig)(conf).Init()
}

func (conf *DatabaseConfig) Build() *gorm.DB {
	url := "/data/db/sqlite/" + conf.Database + ".db"
	if runtime.GOOS == "windows" {
		url = ".." + url
	}
	return (*pkdb.DatabaseConfig)(conf).Build(sqlite.Open(url))
}

type DB pkdb.DB

func (db *DB) Config() any {
	return (*DatabaseConfig)(&db.Conf)
}

func (db *DB) SetEntity(entity interface{}) {
	db.DB = (*DatabaseConfig)(&db.Conf).Build()
}

func (db *DB) Close() error {
	return nil
}
