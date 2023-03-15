package db

import (
	"fmt"
	pkdb "github.com/liov/hoper/server/go/lib/initialize/gormdb"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"runtime"
)

type DatabaseConfig pkdb.DatabaseConfig

func (conf *DatabaseConfig) Build() *gorm.DB {

	var dialector gorm.Dialector

	if conf.Type == pkdb.MYSQL {
		url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
			conf.User, conf.Password, conf.Host,
			conf.Port, conf.Database, conf.Charset)
		dialector = mysql.Open(url)
	} else if conf.Type == pkdb.POSTGRES {
		url := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=%s",
			conf.Host, conf.User, conf.Database, conf.Password, conf.TimeZone)
		dialector = postgres.Open(url)
	} else if conf.Type == pkdb.SQLite {
		url := "/data/db/sqlite/" + conf.Database + ".db"
		if runtime.GOOS == "windows" {
			url = ".." + url
		}
		dialector = sqlite.Open(url)
	}

	return (*pkdb.DatabaseConfig)(conf).Build(dialector)
}

type DB pkdb.DB

func (db *DB) Config() any {
	return (*DatabaseConfig)(&db.Conf)
}

func (db *DB) SetEntity() {
	db.DB = (*DatabaseConfig)(&db.Conf).Build()
}

func (db *DB) Close() error {
	return nil
}
