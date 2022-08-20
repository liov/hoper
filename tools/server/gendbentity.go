package main

import (
	"github.com/actliboy/hoper/server/go/lib/initialize"
	"github.com/actliboy/hoper/server/go/lib/initialize/db/mysql"
	"github.com/actliboy/hoper/server/go/lib/utils/dao/db/gorm/dbtoentity"
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
	dbtoentity.MysqlConvert(Dao.DB.DB)
}
