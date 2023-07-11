package main

import (
	"github.com/hopeio/dora/initialize"
	"github.com/hopeio/dora/initialize/gormdb/mysql"
	"gorm.io/gorm"
)

type Config struct {
}

func (c *Config) Init() {

}

type Dao struct {
	DB mysql.DB `init:"config:MysqlTest"`
}

func (d *Dao) Init() {
	db = d.DB.DB
}

var dao Dao
var db *gorm.DB

func main() {
	defer initialize.Start(nil, &dao)()
	Scan()
	RawScan()
}
