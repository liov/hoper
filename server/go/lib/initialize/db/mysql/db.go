package mysql

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	pkdb "github.com/actliboy/hoper/server/go/lib/initialize/db"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type DatabaseConfig pkdb.DatabaseConfig

func (conf *DatabaseConfig) Generate() any {
	return conf.generate()
}

func (conf *DatabaseConfig) Init() {
	(*pkdb.DatabaseConfig)(conf).Init()
}

func (conf *DatabaseConfig) generate() *gorm.DB {
	url := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		conf.User, conf.Password, conf.Host,
		conf.Port, conf.Database, conf.Charset)
	return (*pkdb.DatabaseConfig)(conf).Generate(mysql.Open(url))
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
