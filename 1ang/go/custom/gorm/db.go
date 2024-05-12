package main

import (
	"github.com/hopeio/cherry/initialize"
	"github.com/hopeio/cherry/initialize/conf_dao/gormdb/mysql"
	"gorm.io/gorm"
)

type Config struct {
	initialize.EmbeddedPresets
}

type Dao struct {
	initialize.EmbeddedPresets
	DB mysql.DB `init:"config:MysqlTest"`
}

func (d *Dao) InitAfterInject() {
	db = d.DB.DB
}

var dao Dao
var db *gorm.DB

func main() {
	defer initialize.Start(&Config{}, &dao)()

	Scan()
	RawScan()
}
