package postgres

import (
	"fmt"
	"github.com/actliboy/hoper/server/go/lib/initialize"
	pkdb "github.com/actliboy/hoper/server/go/lib/initialize/db"
	"gorm.io/driver/postgres"
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
	url := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s TimeZone=%s",
		conf.Host, conf.User, conf.Database, conf.Password, conf.TimeZone)
	return (*pkdb.DatabaseConfig)(conf).Generate(postgres.Open(url))
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
