package main

import (
	"github.com/hopeio/pandora/initialize"
	"github.com/hopeio/pandora/initialize/gormdb/mysql"
	"github.com/hopeio/pandora/utils/dao/db/gorm/mysql/dbtoentity"
)

type config struct {
}

func (c *config) Init() {

}

type dao struct {
	DB mysql.DB `init:"config:MysqlTest""`
}

func (c *dao) Init() {

}

var Dao = &dao{}
var Config = &config{}

func main() {
	defer initialize.Start(Config, Dao)()
	dbtoentity.MysqlConvert(Dao.DB.DB, "./build/generate.go")
}
