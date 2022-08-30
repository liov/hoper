package gorm

import "github.com/actliboy/hoper/server/go/lib/initialize/db/mysql"

type Config struct {
}

func (c *Config) Init() {

}

type Dao struct {
	DB mysql.DB
}

func (d *Dao) Init() {

}

var dao Dao
