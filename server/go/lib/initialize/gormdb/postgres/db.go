package postgres

import (
	"fmt"
	pkdb "github.com/liov/hoper/server/go/lib/initialize/gormdb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatabaseConfig pkdb.DatabaseConfig

func (conf *DatabaseConfig) Init() {
	(*pkdb.DatabaseConfig)(conf).Init()
}

func (conf *DatabaseConfig) Build() *gorm.DB {
	url := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=disable password=%s timezone=%s",
		conf.Host, conf.User, conf.Database, conf.Password, conf.TimeZone)
	return (*pkdb.DatabaseConfig)(conf).Build(postgres.Open(url))
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

func (db *DB) Table(name string) *gorm.DB {
	name = db.Conf.Schema + name
	gdb := db.DB.Clauses()
	gdb.Statement.TableExpr = &clause.Expr{SQL: gdb.Statement.Quote(name)}
	gdb.Statement.Table = name
	return gdb
}

func (db *DB) TableName(name string) string {
	return db.Conf.Schema + name
}

// TODO:
func (db *DB) Inject(configName string) {

}
